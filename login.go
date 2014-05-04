package gamehack

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type User struct {
	AuthorizedWithMoves bool
	MovesUserId int64
	AccessToken string
	RefreshToken string
	Name string `datastore:"-"`
	PhoneEntries []PhoneEntry `datastore:"-"`
}

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		// show sign in screen
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	// if this user isn't in database, add it
	var user User
	err := datastore.Get(c, datastore.NewKey(c, "User", u.ID, 0, nil), &user)
	if err == datastore.ErrNoSuchEntity {
		_, err = datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &User{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
