package gamehack

import (
	"net/http"

	"code.google.com/p/goauth2/oauth"

	"appengine"
	"appengine/urlfetch"
)

var oauthCfg = &oauth.Config{
	ClientId:     clientId,
	ClientSecret: clientSecret,
	AuthURL:      "https://api.moves-app.com/oauth/v1/authorize",
	TokenURL:     "https://api.moves-app.com/oauth/v1/access_token",
	RedirectURL:  "http://localhost:8080/oauth2callback",
	Scope:        "location",
}

func authorize(w http.ResponseWriter, r *http.Request) {
	url := oauthCfg.AuthCodeURL("")
	http.Redirect(w, r, url, http.StatusFound)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	code := r.FormValue("code")

	t := &oauth.Transport{
		Config: oauthCfg,
		Transport: &urlfetch.Transport{
			Context:                       c,
			Deadline:                      0,
			AllowInvalidServerCertificate: false,
		},
	}

	token, err := t.Exchange(code)
	if err != nil {
		c.Errorf(err.Error())
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	tokenCache := cache{
		Context: c,
		Key:     "Oauth",
	}

	err = tokenCache.PutToken(token)
	if err != nil {
		c.Errorf(err.Error())
	}

	t.Token = token

	w.Write([]byte("Authorization flow complete."))
}
