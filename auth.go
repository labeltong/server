package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
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
)

func AuthInitSubrouter(r *mux.Router)  {
	ret := r.PathPrefix("/auth").Subrouter()

	ret.HandleFunc("/login",login)
	ret.HandleFunc("/logout", logout)
	ret.HandleFunc("/secret", secret)


}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
