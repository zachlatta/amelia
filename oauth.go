package amelia

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func authorize(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return nil
	}

	url := oauthCfg.AuthCodeURL(u.ID)
	http.Redirect(w, r, url, http.StatusSeeOther)

	return nil
}

func oauthCallback(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)

	httpErr := r.FormValue("error")
	code := r.FormValue("code")
	userId := r.FormValue("state")

	if httpErr == "access_denied" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	t := CreateTransport(c, nil)

	user := User{}
	userKey := datastore.NewKey(c, "User", userId, 0, nil)

	err := datastore.Get(c, userKey, &user)
	if err != nil {
		return &appError{err, "User not found", http.StatusNotFound}
	}

	token, err := t.Exchange(code)
	if err != nil {
		return &appError{err, "Error authorizing with Moves",
			http.StatusInternalServerError}
	}

	user.MovesToken = *NewMonkeyToken(token)
	user.AuthorizedWithMoves = true

	profile, err := GetUserProfile(t)
	if err != nil {
		return &appError{err, "Error getting user profile from Moves",
			http.StatusInternalServerError}
	}

	user.MovesUserId = profile.UserId

	_, err = datastore.Put(c, userKey, &user)
	if err != nil {
		return &appError{err, "Error saving user", http.StatusInternalServerError}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func revoke(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return nil
	}

	var user User
	userKey := datastore.NewKey(c, "User", u.ID, 0, nil)
	err := datastore.Get(c, userKey, &user)
	if err != nil {
		return &appError{err, "User not found", http.StatusNotFound}
	}

	user.AuthorizedWithMoves = false

	// TODO: See if it's possible to actually revoke access and refresh tokens
	_, err = datastore.Put(c, userKey, &user)
	if err != nil {
		return &appError{err, "Error revoking Moves access", http.StatusInternalServerError}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
