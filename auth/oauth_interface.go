package auth

import (
	"net/http"
)

type OAuthlogin interface{
	HandleMain(w http.ResponseWriter, r *http.Request)
	HandleLoginSuccess(w http.ResponseWriter, r *http.Request)
	HandleLoginFail(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
}



