package main

import (
	"fmt"
	"net/http"
	oauth "./auth"
)
func main() {
	github_oauth_var := oauth.NewGitHubOAuth()
	http.HandleFunc("/", github_oauth_var.HandleMain)
	http.HandleFunc("/login", github_oauth_var.HandleLogin)
	http.HandleFunc("/oauth_callback", github_oauth_var.HandleCallback)
	http.HandleFunc("/login_success", github_oauth_var.HandleLoginSuccess)
	http.HandleFunc("/login_fail", github_oauth_var.HandleLoginFail)
	fmt.Print("Started running on http://127.0.0.1:7000\n")
	fmt.Println(http.ListenAndServe(":7000", nil))
}