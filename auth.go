package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)



/*
/auth/login [POST] Login function, Oauth from client required <- TO BE CONTINUE IF OAUTH_CLIENT IS IMPLEMENTED
/auth/logout [POST] Logout function, Oauth from client required<- TO BE CONTINUE IF OAUTH_CLIENT IS IMPLEMENTED
/auth/secret [GET] Check if user is authenticated
*/

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	secretkey = []byte("super-secret-key")
	store = sessions.NewCookieStore(secretkey)
	authclient *mongo.Client
)

type Credentials struct {
	Email      string   `json:"email"`
	Token      string   `json:"token"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Email      string   `json:"email"`
	jwt.StandardClaims
}

func init(){
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()


	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		os.Getenv("DBadminID"),
		os.Getenv("DBadminPW"),
		os.Getenv("DBHOST"),
		os.Getenv("DBDatabase"),
	)

	authclient, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err !=nil{
		log.Error(fmt.Errorf("error occur in init of dataset handler\n"))
		log.Fatal(err.Error())
	}
	err = client.Connect(ctx)
	if err != nil {
		// FATAL Error : Fail to connect DB
		fmt.Printf("FATAL Error :  mongo client couldn't connect with background context %s", err.Error())
		log.Fatal(err.Error())
	}


}

func AuthInitSubrouter(r *mux.Router)  {
	ret := r.PathPrefix("/auth").Subrouter()

	ret.HandleFunc("/signup",signup)
	ret.HandleFunc("/signout",logout)
	ret.HandleFunc("/testauth", logout)
	ret.HandleFunc("/refresh", secret)


}


func signup(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var result clientdata


	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := client.Database(os.Getenv("DBDatabase")).Collection("client_list")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res := collection.FindOne(ctx, bson.D{{"email" , creds.Email}})

	err = res.Decode(&result)

	if err != nil{
		log.Error(fmt.Errorf("error occur Decoding found datatolabel struct\n"))
		log.Fatal(err.Error())
	}

	if result.Email != "" {
		// If already signed,  return 200 OK and "Already exist email." on body.
		w.WriteHeader(http.StatusOK)
		_, _  = w.Write([]byte("Already exist email."))
		return
	}else{
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Email: creds.Email,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})



	}




	// Authentication goes here
	// ...

	// Set user as authenticated

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
