package amelia

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type User struct {
	AuthorizedWithMoves bool
	MovesToken          MonkeyToken
	MovesUserId         int64
	LastAddress         string
	Name                string       `datastore:"-"`
	PhoneEntries        []PhoneEntry `datastore:"-"`
}

func login(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		// show sign in screen
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return nil
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	// if this user isn't in database, add it
	var user User
	err := datastore.Get(c, datastore.NewKey(c, "User", u.ID, 0, nil), &user)
	if err == datastore.ErrNoSuchEntity {
		_, err = datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &User{})
		if err != nil {
			return &appError{err, "Error creating user", http.StatusInternalServerError}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func logout(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		url, err := user.LogoutURL(c, "/")
		if err != nil {
			return &appError{err, "Error logging out", http.StatusInternalServerError}
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
