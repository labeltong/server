package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"

	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	coll *mongo.Collection
	client *mongo.Client

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey

	privKeyPath = "secrets/labeltong.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "secrets/labeltong.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub

)

func AuthInitSubrouter(r *mux.Router)  {
	ret := r.PathPrefix("/auth").Subrouter()

	ret.HandleFunc("/login",LoginHandler).Methods("POST")
	ret.HandleFunc("/testauth", TestAuthHandler).Methods("GET")
	//ret.HandleFunc("/profile", ProfileAuthHandler).Methods("GET")


}

func init()  {
	var err error
	client, err = ConnectDBClient()
	if err != nil{
		log.Fatal(err.Error())
	}
	coll = client.Database(os.Getenv("DBDatabase")).Collection("client_list")

	// For public/private key

	signBytes, err := ioutil.ReadFile(privKeyPath)
	customfatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	customfatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	customfatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	customfatal(err)

}



func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user user
	var result clientdata
	var res ResponseResult

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	if err != nil {
		res.Error = err.Error()
		err = json.NewEncoder(w).Encode(res)
		if err != nil{
			log.Fatal(err.Error())

		}
		return
	}


	err = coll.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if user.Email == result.Email{ // If email already exists -> refresh token
		w.WriteHeader(http.StatusOK)
		//_,err = w.Write([]byte("User already exists"))
		//
		//if err != nil {
		//	res.Error = err.Error()
		//	err = json.NewEncoder(w).Encode(res)
		//	if err != nil{
		//		log.Fatal(err.Error())
		//
		//	}
		//	return
		//}
		newclient:= clientdata{ //Make new user's data
			Email:user.Email,
			Token:user.Token,
			Points:0,
			Banpoint:0,
			Isbanned:false,
			IsAdmin:false,
			Pointusage:[]string{},
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Email: newclient.Email,
			Points:newclient.Points,
			Isbanned:newclient.Isbanned,
			Pointusage:newclient.Pointusage,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
		// Create the JWT string
		tokenString, err := token.SignedString(signKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintln(w, "Sorry, error while Signing Token!")
			customfatal(err)
			log.Printf("Token Signing error: %v\n", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(tokenString))

	}else{ // Register new users
		newclient:= clientdata{ //Make new user's data
			Email:user.Email,
			Token:user.Token,
			Points:0,
			Banpoint:0,
			Isbanned:false,
			IsAdmin:false,
			Pointusage:[]string{},
		}


		_,err = coll.InsertOne(context.TODO(),newclient) //Insert new user's data to db -> client_list

		if err != nil {
			res.Error = err.Error()
			err = json.NewEncoder(w).Encode(res)
			if err != nil{
				log.Fatal(err.Error())

			}
			return
		}
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Email: newclient.Email,
			Points:newclient.Points,
			Isbanned:newclient.Isbanned,
			Pointusage:newclient.Pointusage,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
		// Create the JWT string
		tokenString, err := token.SignedString(signKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintln(w, "Sorry, error while Signing Token!")
			customfatal(err)
			log.Printf("Token Signing error: %v\n", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(tokenString))
		customfatal(err)
	}

}

func ValidateToken(w http.ResponseWriter, r *http.Request)  error {


		tokenString := r.Header.Get("Authorization")
		//1. No token
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_,err := w.Write([]byte("Missing Authorization Header"))
			customfatal(err)
		}
		//2. Token exist
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return verifyKey, nil
		})
		switch err.(type){
		case nil:
			if token.Valid{
				//2.1 valid token
				w.WriteHeader(http.StatusOK)
				return nil

			}else{
				//2.2 invalid token
				w.WriteHeader(http.StatusUnauthorized)
				_, e := fmt.Fprintln(w, "Token not valid.")
				customfatal(e)
				return err


			}

		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				//2.3 expired token

				w.WriteHeader(http.StatusUnauthorized)
				_, e := fmt.Fprintln(w, "Token Expired, get a new one.")
				customfatal(e)
				return  err


			default:
				w.WriteHeader(http.StatusInternalServerError)
				_,e := fmt.Fprintln(w, "Error while Parsing Token!")
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				customfatal(e)
				return err

			}
		default: // something else went wrong
			w.WriteHeader(http.StatusInternalServerError)
			_,e := fmt.Fprintln(w, "Error while Parsing Token!")
			customfatal(e)
			log.Printf("Token parse error: %v\n", err)
			return err
		}


		//name := token.Claims.(jwt.MapClaims)

}

func TestAuthHandler(w http.ResponseWriter, r *http.Request) {
	err := ValidateToken(w,r)
	if err !=nil{
		_,_ = w.Write([]byte("forbidden") )


	}else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("authed"))
	}
}

//func ProfileAuthHandler(w http.ResponseWriter, r *http.Request) {
//	t, err := ValidateToken(w,r)
//	if err !=nil{
//		_,_ = w.Write([]byte("forbidden") )
//		return
//
//	}else {
//		_,_ = w.Write([]byte("forbidden") )
//
//		}
//	}
//}
//func ProfileHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	tokenString := r.Header.Get("Authorization")
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		// Don't forget to validate the alg is what you expect:
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method")
//		}
//		return []byte("secret"), nil
//	})
//	var result model.User
//	var res model.ResponseResult
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		result.Username = claims["username"].(string)
//		result.FirstName = claims["firstname"].(string)
//		result.LastName = claims["lastname"].(string)
//
//		json.NewEncoder(w).Encode(result)
//		return
//	} else {
//		res.Error = err.Error()
//		json.NewEncoder(w).Encode(res)
//		return
//	}
//
//}