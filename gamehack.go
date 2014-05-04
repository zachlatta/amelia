package gamehack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/subosito/twilio"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"time"
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

type Segment struct {
	Place      Place  `json:"place"`
	LastUpdate string `json:"lastUpdate"`
}

type DailySegments struct {
	Segments []Segment `json:"segments"`
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/phone", phone)
	http.HandleFunc("/addphone", addPhone)
	http.HandleFunc("/delphone", delPhone)
	http.HandleFunc("/authorize", authorize)
	http.HandleFunc("/oauth2callback", oauthCallback)
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

// TODO: this function is untested
func updateDailySegments(dailySegmentsList []DailySegments, userID string, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// get last update time from database
	var user User
	userKey := datastore.NewKey(c, "User", userID, 0, nil)
	err := datastore.Get(c, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastUpdate, err := time.Parse(time.RFC3339Nano, user.LastUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get phone numbers from database
	var phoneEntries []PhoneEntry
	_, err = datastore.NewQuery("PhoneEntry").Ancestor(datastore.NewKey(c, "User", userID, 0, nil)).GetAll(c, &phoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, dailySegments := range dailySegmentsList {
		for _, segment := range dailySegments.Segments {
			// send texts to user's phone numbers
			for _, phone := range phoneEntries {
				sendText(segment.Place, phone.Phone, w, r)
			}
			// update last update time
			time, err := time.Parse(time.RFC3339Nano, segment.LastUpdate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if time.After(lastUpdate) {
				lastUpdate = time
			}
		}
	}
	// set last update time in database
	// TODO: this introduces a race condition setting LastUpdate if updateDailySegments() is called in close succession for the same user (may want to hold lock on database)
	user.LastUpdate = lastUpdate.Format(time.RFC3339Nano)
	_, err = datastore.Put(c, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
