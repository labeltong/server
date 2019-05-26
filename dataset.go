package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	datadbs postgredb

	dataset_verifyKey *rsa.PublicKey
	//dataset_signKey   *rsa.PrivateKey
)

func init() {
	//
	//signBytes, err := ioutil.ReadFile(privKeyPath)
	//Custom_fatal(err)
	//
	//asignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	//Custom_fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	Custom_fatal(err)

	dataset_verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	Custom_fatal(err)

}

//routes
// 2. ListTags /list/tags [GET] -> list of tags by JSON
// 3. GetByMethods /list/methods/:id [GET] -> get data by methods
// 4. GetByTags /list/tags/:id [GET] -> get data by tags
// optional
// 5. Listdataset /list/dataset[GET] -> list of datasets by JSON
// 6. GetBydatasets /list/methods/:id [GET] -> get data by methods

//TODO : MAKE example column to db and model
//TODO : Data to S3 bucket and modify value

//answer type
//1 box 2 classification 3 sentiment(audio)

func ListTags(c echo.Context) error {
	var results []TagList
	tmpacc := datadbs.DB.Find(&results)
	if tmpacc.Error != nil {
		fmt.Printf("Time : %s [500 Error] Error while search all tags \n", time.Now().String())
		return echo.ErrInternalServerError
	}
	pagesJson, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("Time : %s [500 Error] Error while marshal all tags data to JSON \n", time.Now().String())
		return echo.ErrInternalServerError
	}

	return c.JSONBlob(http.StatusOK, pagesJson)
}

//answer type
//1 box 2 classification 3 sentiment(audio)
func GetByMethods(c echo.Context) error {
	var res Datas
	var isdup int
	Methods := c.Param("method")
	usertoken := strings.Split(c.Request().Header["Authorization"][0], " ")[1]
	claims := UserInfoClaim{}
	tknstr, _ := jwt.ParseWithClaims(usertoken, &claims, func(token *jwt.Token) (interface{}, error) {
		return dataset_verifyKey, nil
	})
	if !tknstr.Valid {
		return echo.ErrUnauthorized
	}
	fmt.Println(usertoken)
	user := claims.Email
	emailtoid := User{}
	err := datadbs.DB.Select("id").Where("email = ?", user).First(&emailtoid)
	Custom_panic(err.Error)

	datadbs.DB.Where("answer_type = ?", Methods).Where("required_num_answer > ?", 0).Order(gorm.Expr("random()")).First(&res)
	//TODO: NO dups data for user
	datadbs.DB.Where("id = ?", emailtoid.ID).Where("data_id = ?", res.ID).Find(&AllAnswers{}).Count(&isdup)
	if isdup != 0 {
		fmt.Printf("Time : %s [500 Error] Already seen this data \n", time.Now().String())
		return echo.ErrInternalServerError
	}
	fmt.Println(res)
	return c.JSON(http.StatusOK, res)

}

func GetByTags(c echo.Context) error {
	var res Datas
	var tagres Tags
	var isdup int
	tags := c.Param("tag")
	usertoken := strings.Split(c.Request().Header["Authorization"][0], " ")[1]
	claims := UserInfoClaim{}
	tknstr, _ := jwt.ParseWithClaims(usertoken, &claims, func(token *jwt.Token) (interface{}, error) {
		return dataset_verifyKey, nil
	})
	if !tknstr.Valid {
		return echo.ErrUnauthorized
	}
	fmt.Println(usertoken)
	user := claims.Email
	emailtoid := User{}
	err := datadbs.DB.Select("id").Where("email = ?", user).First(&emailtoid)
	Custom_panic(err.Error)

	datadbs.DB.Where("tag_id = ?", tags).Order(gorm.Expr("random()")).First(&tagres)
	i := tagres.DataID
	datadbs.DB.Where("id = ?", i).Where("required_num_answer > ?", 0).Order(gorm.Expr("random()")).First(&res)
	//TODO: NO dups data for user

	datadbs.DB.Where("id = ?", emailtoid.ID).Where("data_id = ?", res.ID).Find(&AllAnswers{}).Count(&isdup)
	if isdup != 0 {
		fmt.Printf("Time : %s [500 Error] Already seen this data \n", time.Now().String())
		return echo.ErrInternalServerError
	}
	fmt.Println(res)
	return c.JSON(http.StatusOK, res)

}

//
//
func main() {
	e := echo.New()
	port := os.Args[1]

	datadbs = postgredb{}
	err := datadbs.Connect()
	Custom_fatal(err)
	defer func() {
		internalerr := datadbs.DB.Close()
		panic(internalerr)

	}()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"localhost"},
	//	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	//	AllowHeaders:[]string{"*"},
	//	ExposeHeaders:[]string{"*"},
	//}))

	r := e.Group("/dataset")
	r.GET("/tags", ListTags)
	r.GET("/dmethods/:method", GetByMethods)
	r.GET("/dtags/:tag", GetByTags)
	//Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:        &UserInfoClaim{},
		SigningKey:    dataset_verifyKey,
		SigningMethod: "RS256",
	}
	r.Use(middleware.JWTWithConfig(config))

	e.Logger.Fatal(e.Start(port))

}
