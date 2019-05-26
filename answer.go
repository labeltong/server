package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	answerdbs postgredb

	answer_verifyKey *rsa.PublicKey
	//dataset_signKey   *rsa.PrivateKey
)

func init() {

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	Custom_fatal(err)

	answer_verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	Custom_fatal(err)

}

func CreateAnswer(db *gorm.DB, ans AnswersJSON) error {
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var emailtoid User
	err := db.Where("email = ?", ans.Email).First(&emailtoid)
	Custom_panic(err.Error)

	q := Answers{
		UserId:     emailtoid.ID,
		DataId:     ans.DataId,
		AnswerData: ans.AnswerData,
		AnswerTime: time.Now(),
	}

	if err := tx.Create(&q).Error; err != nil {
		tx.Rollback()
		return err
	}
	if a := tx.Table("public.datas").Where("id = ?", q.DataId).UpdateColumn("required_num_answer", gorm.Expr("required_num_answer - ?", 1)).Error; a != nil {
		tx.Rollback()
		return a
	}

	if b := tx.Table("public.users").Where("id = ?", q.UserId).UpdateColumn("points", gorm.Expr("points + ?", 1)).Error; b != nil {
		tx.Rollback()
		return b
	}

	return tx.Commit().Error

}

func SubmitData(c echo.Context) error {
	//TODO: Put data to DB
	//TODO
	submited := AnswersJSON{}

	if err := c.Bind(&submited); err != nil {
		fmt.Println(submited)
		fmt.Printf("Time : %s [500 Error] Error while intepreting request's json data \n", time.Now().String())
		return echo.ErrInternalServerError
	}
	//ansdb := answerdbs.DB.Table("public.answer_table")
	ansdb := answerdbs.DB

	err := CreateAnswer(ansdb, submited)

	if err != nil {
		fmt.Printf("Time : %s [500 Error] Error while posting request's json data \n", time.Now().String())
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, submited)

}

func main() {

	e := echo.New()
	port := os.Args[1]

	answerdbs = postgredb{}
	err := answerdbs.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		internalerr := answerdbs.DB.Close()
		if internalerr != nil {
			log.Panic(internalerr)

		}

	}()

	// Middleware
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"localhost"},
	//	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	//	AllowHeaders:[]string{"*"},
	//	ExposeHeaders:[]string{"*"},
	//}))

	//Restricted group
	r := e.Group("/answer")
	r.POST("/answer", SubmitData)

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:        &UserInfoClaim{},
		SigningKey:    answer_verifyKey,
		SigningMethod: "RS256",
	}
	r.Use(middleware.JWTWithConfig(config))
	//r.GET("", restricted)

	e.Logger.Fatal(e.Start(port))

}
