package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// login -> Get token && refresh token  [GET] tkn x
// register -> Make token and return [POST] tkn x
// unregister -> delete member [DELETE] tkn 0
// edit -> edit member's info, [PUT] tkn 0
//refactortest
var (
	authdbs postgredb
	//upgrader = websocket.Upgrader{}

	averifyKey *rsa.PublicKey
	asignKey   *rsa.PrivateKey
)

func init() {

	signBytes, err := ioutil.ReadFile(privKeyPath)
	Custom_fatal(err)

	asignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	Custom_fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	Custom_fatal(err)

	averifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	Custom_fatal(err)

}

func LoginHandler(c echo.Context) error {
	// get, email=email, token=token,
	// login and refresh

	email, token := c.QueryParam("email"), c.QueryParam("token")
	//email, token = fmt.Sprintf(`"` + email + `"`),fmt.Sprintf(`"` + token +`"`)

	///HARDCODED_SHIT REMOVE ASAP!!!
	//TODO: Remove this code asap

	if token == "ksh" {

	}

	///HARDCODED_SHIT REMOVE ASAP!!!

	//ACCESS DATABASE
	//TODO: Database Access to common db instance
	res := User{}
	tmpacc := authdbs.DB.Where("email = ?", email).First(&res)
	ifnotexist := tmpacc.RecordNotFound()
	fmt.Println(tmpacc.Value)

	if ifnotexist {
		fmt.Printf("Time : %s [404 Error] Email not found email=%s\n", time.Now().String(), email)
		return echo.ErrNotFound
	}

	v := googleAuthIDTokenVerifier.Verifier{}
	fmt.Println(token)
	err := v.VerifyIDToken(token, []string{})
	if err == googleAuthIDTokenVerifier.ErrInvalidToken {
		fmt.Println("Invalid token")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", email))
	} else if err == googleAuthIDTokenVerifier.ErrWrongSignature {
		fmt.Println("wrong signature")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", email))

	} else if err == googleAuthIDTokenVerifier.ErrTokenUsedTooLate {
		fmt.Println("token expired")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", email))
	} else if err != nil {
		log.Error(err.Error())
	}
	claimSet, _ := googleAuthIDTokenVerifier.Decode(token)

	if claimSet.Email != email {
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", email))

	}

	//ACCESS DATABASE

	fmt.Printf("Time : %s [200 OK] autorized login with email=%s, token=%s\n", time.Now().String(), email, token)
	claims := &UserInfoClaim{

		Email:         res.Email,
		Token:         "ABRACADABRA!",
		Name:          res.Name,
		PhoneNum:      res.PhoneNum,
		Points:        res.Points,
		IsBanned:      res.IsBanned,
		RegisterDate:  res.CreatedAt,
		LastLoginDate: res.LastLoginDate,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * expiredTimeInMinute).Unix()
	claims.IssuedAt = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenstring, tknerr := t.SignedString(asignKey)

	if tknerr != nil {
		fmt.Printf("Time : %s [500 Error] Error while signing jwt token with email=%s\n", time.Now().String(), email)
		return echo.ErrInternalServerError

	}
	//ACCESS DATABASE

	loginchange := authdbs.DB.Model(&res).Where("id = ?", res.ID).Update("last_login_date", time.Now())

	//ACCESS DATABASE

	if loginchange.Error != nil {
		fmt.Printf("Time : %s [500 Error] Error while changing last_login date with email=%s\n", time.Now().String(), email)
		return echo.ErrInternalServerError

	}

	fmt.Printf("Time : %s [200 OK] Token publish successfully with user email=%s\n", time.Now().String(), email)
	return c.String(http.StatusOK, tokenstring)

}

func RegisterHandler(c echo.Context) error {

	newone := &MakeUser{}
	if err := c.Bind(newone); err != nil {
		fmt.Printf("Time : %s [500 Error] Error while creating user  with email=%s, token=%s\n", time.Now().String(), newone.Email, newone.Token)
		return echo.ErrInternalServerError
	}

	tmpdb := authdbs.DB.Where("email = ?", newone.Email).First(&User{})

	if tmpdb.Value == nil {
		fmt.Printf("Time : %s [500 Error] User exists  with email=%s \n", time.Now().String(), newone.Email)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Email %s Already exist", newone.Email))
	}

	v := googleAuthIDTokenVerifier.Verifier{}
	err := v.VerifyIDToken(newone.Token, []string{})
	if err == googleAuthIDTokenVerifier.ErrInvalidToken {
		fmt.Println("Invalid token")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", newone.Email))
	} else if err == googleAuthIDTokenVerifier.ErrWrongSignature {
		fmt.Println("wrong signature")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", newone.Email))

	} else if err == googleAuthIDTokenVerifier.ErrTokenUsedTooLate {
		fmt.Println("token expired")
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", newone.Email))
	} else if err != nil {
		log.Error(err.Error())
	}
	claimSet, _ := googleAuthIDTokenVerifier.Decode(newone.Token)

	if claimSet.Email != newone.Email {
		return c.String(echo.ErrForbidden.Code, fmt.Sprintf("Email %s is not valid with token.", newone.Email))

	}

	tmpacc := authdbs.DB.Create(&User{
		Email:    newone.Email,
		Token:    "ABRACADABRA!",
		Name:     newone.Name,
		PhoneNum: newone.PhoneNum,

		LastLoginDate: time.Now(),
	})

	if tmpacc.Error != nil {
		fmt.Printf("Time : %s [500 Error] Error while creating user  with email=%s, token=%s\n", time.Now().String(), newone.Email, newone.Token)
		return echo.ErrInternalServerError
		//fatal(tmpacc.Error)

	}

	fmt.Printf("Time : %s [201 created] Creating user  with email=%s, \n", time.Now().String(), newone.Email)
	return c.JSON(http.StatusCreated, newone)

}

func main() {
	//name := os.Args[1]
	port := os.Args[1]
	e := echo.New()
	authdbs = postgredb{}
	err := authdbs.Connect()
	Custom_panic(err)
	defer func() {
		internalerr := authdbs.DB.Close()
		Custom_panic(internalerr)

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

	// Login route
	e.GET("/auth", LoginHandler)
	e.POST("/auth", RegisterHandler)

	e.Logger.Fatal(e.Start(port))
}
