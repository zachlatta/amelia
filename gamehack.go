package gamehack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/subosito/twilio"

	"appengine"
	"appengine/urlfetch"
)

type Notification struct {
	UserID           int64             `json:"userId"`
	StorylineUpdates []StorylineUpdate `json:"storylineUpdates"`
}

type StorylineUpdate struct {
	// TODO: Change to equivalent of enum
	Reason          string `json:"reason"`
	LastSegmentType string `json:"lastSegmentType"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Place struct {
	Id       int      `json:"id"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

var oauthCfg = &oauth.Config{
	ClientId:     clientId,
	ClientSecret: clientSecret,
	AuthURL:      "https://api.moves-app.com/oauth/v1/authorize",
	TokenURL:     "https://api.moves-app.com/oauth/v1/access_token",
	RedirectURL:  "http://localhost:8080/oauth2callback",
	Scope:        "location",
}

func init() {
	http.HandleFunc("/authorize", authorize)
	http.HandleFunc("/oauth2callback", oauthCallback)
	http.HandleFunc("/notification", handleNotification)
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

func handleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body.", http.StatusBadRequest)
		return
	}

	var notification Notification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hasDataUpload := false
	for _, update := range notification.StorylineUpdates {
		if update.Reason == "DataUpload" {
			hasDataUpload = true
			break
		}
	}

	if hasDataUpload {

	}
	/*fmt.Fprintf(w, "%v", notification)
	if err != nil {
		http.Error(w, "Error writing response body.", http.StatusInternalServerError)
		return
	}*/
	/*var place Place
	err = json.Unmarshal(body, &place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendText(place, "+15555555555", w, r)*/
}

func sendText(place Place, phone string, w http.ResponseWriter, r *http.Request) {
	a := appengine.NewContext(r)
	f := urlfetch.Client(a)
	c := twilio.NewClient(twilioSid, twilioAuthToken, f)

	params := twilio.MessageParams{
		Body: fmt.Sprintf("Your child is now at lat %f lon %f", place.Location.Lat, place.Location.Lon),
	}
	_, _, err := c.Messages.Send("+15555555555", phone, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
