package gamehack

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/subosito/twilio"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"appengine/user"
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

type UserInfo struct {
	User         string
	PhoneEntries []PhoneEntry
}

type PhoneEntry struct {
	Parent string
	Phone  string
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

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, homePage)
}

const homePage = `
<!doctype html>
<html>
  <head>
    <title>Game+Hack</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
  </head>
  <body>
    <form action="/login" method="post">
      <div><input type="submit" value="Sign In"></div>
    </form>
  </body>
</html>
`

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	http.Redirect(w, r, "/phone", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		url, err := user.LogoutURL(c, "/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

var phoneTemplate = template.Must(template.New("phone").Parse(`
<!doctype html>
<html>
  <head>
    <title>Game+Hack</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
  </head>
  <body>
    <p>Hello, {{.User}}! <a href="/logout">Sign Out</a></p>
    <form action="/addphone" method="POST">
      <div>Parent: <input type="text" name="parent"/></div>
      <div>Phone: <input type="text" name="phone"/></div>
      <div><input type="submit" value="Add Phone Number"></div>
    </form>
    <form action="/delphone" method="POST">
      <select name="parent">
        <option value=""></option>
        {{range .PhoneEntries}}
          <option value="{{.Parent}}">{{.Parent}}</option>
        {{end}}
      </select>
      <div><input type="submit" value="Remove Phone Number"></div>
    </form>
    {{range .PhoneEntries}}
      <p><b>{{.Parent}}</b>: {{.Phone}}</p>
    {{end}}
  </body>
</html>
`))

func phone(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	userInfo := UserInfo{
		User: u.Email,
	}
	_, err := datastore.NewQuery("PhoneEntry").Ancestor(datastore.NewKey(c, "User", u.ID, 0, nil)).GetAll(c, &userInfo.PhoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = phoneTemplate.Execute(w, &userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addPhone(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	_, err := datastore.Put(c, datastore.NewKey(c, "PhoneEntry", r.FormValue("parent"), 0, datastore.NewKey(c, "User", u.ID, 0, nil)), &PhoneEntry{
		Parent: r.FormValue("parent"),
		Phone:  r.FormValue("phone"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/phone", http.StatusFound)
}

func delPhone(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	err := datastore.Delete(c, datastore.NewKey(c, "PhoneEntry", r.FormValue("parent"), 0, datastore.NewKey(c, "User", u.ID, 0, nil)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/phone", http.StatusFound)
}
