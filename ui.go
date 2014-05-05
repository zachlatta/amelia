package gamehack

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Base directory where templates are stored.
const tD = "static/"

var templates = template.Must(template.ParseFiles(
	tD+"index.html",
	tD+"phone.html",
))

func root(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func phone(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	user := User{
		Name: u.String(),
	}

	userKey := datastore.NewKey(c, "User", u.ID, 0, nil)

	err := datastore.Get(c, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = datastore.NewQuery("PhoneEntry").Ancestor(userKey).GetAll(c, &user.PhoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "phone", user)
}

func renderTemplate(w http.ResponseWriter, tmpl string, c interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
