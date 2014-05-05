package gamehack

import (
	"fmt"
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

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
	t, err := template.ParseFiles(fmt.Sprintf("static/%s.html", tmpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
