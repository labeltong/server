package auth

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	_ "golang.org/x/oauth2/google"
	"net/http"
)



type Github_oauth struct {
	oauthConf oauth2.Config
	oauthStateString string
}

func NewGitHubOAuth() *Github_oauth{
	o := Github_oauth{}
	o.oauthConf = oauth2.Config{
		ClientID:     "YOUR CLIENT ID",
		ClientSecret: "YOUR CLIENT SECRET",
		Scopes:       []string{},
		Endpoint:     githuboauth.Endpoint,
	}
	o.oauthStateString = RandToken()
	return &o
}


// /
func (o *Github_oauth)HandleMain(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "static/main.html", "GitHub")
}

// /login_success
func (o *Github_oauth)HandleLoginSuccess(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "static/login_success.html", "GitHub")

}
// /login_fail
func (o *Github_oauth)HandleLoginFail(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "static/login_Fail.html", "GitHub")

}


// /login
func (o *Github_oauth)HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := o.oauthConf.AuthCodeURL(o.oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
// / /oauth_callback
func (o *Github_oauth)HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	if state != o.oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", o.oauthStateString, state)
		http.Redirect(w, r, "/login_fail", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := o.oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/login_fail", http.StatusTemporaryRedirect)
		return
	}
	ctx := context.Background()

	oauthClient := o.oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx,"")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/login_fail", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/login_success", http.StatusTemporaryRedirect)
}
