package amelia

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
	tD+"application.html",
	tD+"landing.html",
	tD+"profile.html",
))

func root(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		return renderTemplate(w, "landing", nil)
	}

	user := User{
		Name: u.String(),
	}

	userKey := datastore.NewKey(c, "User", u.ID, 0, nil)

	err := datastore.Get(c, userKey, &user)
	if err != nil {
		return &appError{err, "User not found", http.StatusNotFound}
	}

	_, err = datastore.NewQuery("PhoneEntry").Ancestor(userKey).GetAll(c, &user.PhoneEntries)
	if err != nil {
		return &appError{err, "Error retrieving phones", http.StatusInternalServerError}
	}
	return renderTemplate(w, "profile", user)
}

func renderTemplate(w http.ResponseWriter, tmpl string,
	c interface{}) *appError {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		return &appError{err, "Can't display webpage", http.StatusInternalServerError}
	}
	return nil
}
