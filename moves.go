package amelia

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"

	"appengine"
	"appengine/urlfetch"
)

const (
	baseURL = "https://api.moves-app.com/api/1.1"
)

var oauthCfg = &oauth.Config{
	ClientId:     clientId,
	ClientSecret: clientSecret,
	AuthURL:      "https://api.moves-app.com/oauth/v1/authorize",
	TokenURL:     "https://api.moves-app.com/oauth/v1/access_token",
	RedirectURL:  "http://localhost:8080/oauth2callback",
	Scope:        "location",
}

type RFC3339Time struct {
	time.Time
}

func (t *RFC3339Time) UnmarshalJSON(b []byte) (err error) {
	t.Time, err = time.Parse("20060102T150405-0700", strings.Trim(string(b), `"`))
	return err
}

func (t *RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"` + time.RFC3339 + `"`)), nil
}

func CreateTransport(c appengine.Context, token *oauth.Token) *oauth.Transport {
	return &oauth.Transport{
		Config: oauthCfg,
		Transport: &urlfetch.Transport{
			Context:                       c,
			Deadline:                      0,
			AllowInvalidServerCertificate: false,
		},
		Token: token,
	}
}

type UserProfile struct {
	UserId int64 `json:"userId"`
}

func GetUserProfile(t *oauth.Transport) (*UserProfile, error) {
	resp, err := t.Client().Get(baseURL + "/user/profile")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	profile := new(UserProfile)
	err = json.Unmarshal(bytes, &profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func GetLatestPlaces(t *oauth.Transport) (*[]DailySegments, error) {
	resp, err := t.Client().Get(baseURL + "/user/places/daily?pastDays=5")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dailySegments []DailySegments
	err = json.Unmarshal(bytes, &dailySegments)
	if err != nil {
		return nil, err
	}

	return &dailySegments, nil
}
