package gamehack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
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
	Id int `json:"id"`
	Type string `json:"type"`
	Location Location `json:"location"`
}

var (
	AccountSid = "AccountSid goes here"
	AuthToken  = "AuthToken goes here"
)

func init() {
	http.HandleFunc("/notification", handleNotification)
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
	// Initialize twilio client
	//c := twilio.NewClient(AccountSid, AuthToken, nil)

	// You can set custom Client, eg: you're using `appengine/urlfetch` on Google's appengine
	a := appengine.NewContext(r) // r is a *http.Request
	f := urlfetch.Client(a)
	c := twilio.NewClient(AccountSid, AuthToken, f)

	// Send Message
	params := twilio.MessageParams{
		Body: fmt.Sprintf("Your child is now at lat %f lon %f", place.Location.Lat, place.Location.Lon),
	}
	_, _, err := c.Messages.Send("+15555555555", phone, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
