
package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)
var (
	dbs postgredb

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)


func init(){

	signBytes, err := ioutil.ReadFile(privKeyPath)
	Custom_fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	Custom_fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	Custom_fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	Custom_fatal(err)
}

func InfoHandler(c echo.Context) error {
	// get, email=email, token=token,
	// login and refresh


	usertoken := strings.Split(c.Request().Header["Authorization"][0]," ")[1]
	claims := UserInfoClaim{}
	tknstr, _ := jwt.ParseWithClaims(usertoken, &claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if !tknstr.Valid {
		return echo.ErrUnauthorized
	}
	fmt.Println(usertoken)
	user :=claims.Email
	emailtoid := User{}
	err := dbs.DB.Where("email = ?",user).First(&emailtoid)
	fmt.Println(emailtoid)
	if err.Error != nil{
		return c.String(http.StatusInternalServerError,"server db failed")
	}
	return c.JSON(http.StatusOK,emailtoid)
}


func main() {
	e := echo.New()
	port := os.Args[1]


	dbs = postgredb{}
	err := dbs.Connect()
	Custom_fatal(err)
	defer func(){
		internalerr := dbs.DB.Close()
		Custom_fatal(internalerr)

	}()

	// Setup middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	//optional middleware

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:[]string{"*"},
		ExposeHeaders:[]string{"*"},
	}))
	//e.Use(middleware.Gzip())
	//e.Use(middleware.CSRF())

	//Set groups

	authg := e.Group("/auth")
	datag := e.Group("/dataset")
	answg := e.Group("/answer")
	e.GET("/info",InfoHandler)
	//for item usage
	//itemg := e.Group("/item")


	// Configure middleware with the custom claims type
	//jwtconfig := middleware.JWTConfig{
	//	Claims:     &UserInfoClaim{},
	//	SigningKey: verifyKey,
	//	SigningMethod:"RS256",
	//}
	//datag.Use(middleware.JWTWithConfig(jwtconfig)
	//answg.Use(middleware.JWTWithConfig(jwtconfig))

	// Setup reverse proxy

	// /auth proxy
	_ ,authProxyTarget := MakeReverseProxy(AuthRP,&NumAuthRP)
	authg.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(authProxyTarget)))

	// /dataset proxy
	_ ,dataProxyTarget := MakeReverseProxy(DataRP,&NumDataRP)
	datag.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(dataProxyTarget)))

	// 	/answer proxy
	_ ,ansProxyTarget := MakeReverseProxy(AnswerRP,&NumAnswerRP)
	answg.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(ansProxyTarget)))


	e.Logger.Fatal(e.Start(port))
}