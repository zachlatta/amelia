package amelia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/subosito/twilio"
	"github.com/zachlatta/go-tomtom"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

type Notification struct {
	UserID           int64             `json:"userId"`
	StorylineUpdates []StorylineUpdate `json:"storylineUpdates"`
}

type StorylineUpdate struct {
	// TODO: Change to equivalent of enum
	Reason          string `json:"reason"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Place struct {
	Location Location `json:"location"`
}

type Segment struct {
	Place      Place  `json:"place"`
}

type DailySegments struct {
	Segments []Segment `json:"segments"`
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

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
		q := datastore.NewQuery("User").Filter("AuthorizedWithMoves =", true).Filter("MovesUserId =", notification.UserID)

		var users []User
		keys, err := q.GetAll(c, &users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(users) <= 0 {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		user, key := users[0], keys[0]

		t := CreateTransport(c, user.MovesToken.NormToken())

		dailySegmentsList, err := GetLatestPlaces(t)
		if err != nil {
			c.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updateDailySegments(*dailySegmentsList, user, key, w, r)
	}
}

func updateDailySegments(dailySegmentsList []DailySegments, user User, userKey *datastore.Key, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	f := urlfetch.Client(c)
	var phoneEntries []PhoneEntry
	_, err := datastore.NewQuery("PhoneEntry").Ancestor(userKey).GetAll(c, &phoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, dailySegments := range dailySegmentsList {
		for _, segment := range dailySegments.Segments {
			// use reverse geocoding to find address
			client := tomtom.NewClient(tomtomKey, f)
			codes, err := client.Geocode.ReverseGeocode(segment.Place.Location.Lat, segment.Place.Location.Lon)
			if err != nil {
				c.Errorf(err.Error())
			}

			// get address string and check if it changed
			var address string
			if err == nil && len(codes) > 0 {
				address = codes[0].FormattedAddress
			} else {
				address = fmt.Sprintf("%f, %f", segment.Place.Location.Lat, segment.Place.Location.Lon)
			}
			if address == user.LastAddress {
				return
			}
			user.LastAddress = address
			_, err = datastore.Put(c, userKey, &user)
			if err != nil {
				c.Errorf(err.Error())
			}

			// send texts
			for _, phone := range phoneEntries {
				sendText("I'm now at " + address + ".", phone.Phone, w, r)
			}
			return // first item in slice is latest one
		}
	}
}

func sendText(message string, phone string, w http.ResponseWriter, r *http.Request) {
	a := appengine.NewContext(r)
	f := urlfetch.Client(a)
	c := twilio.NewClient(twilioSid, twilioAuthToken, f)

	var params twilio.MessageParams
	params.Body = message
	_, _, err := c.Messages.Send(twilioPhone, phone, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
