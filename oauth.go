package gamehack

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func authorize(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	url := oauthCfg.AuthCodeURL(u.ID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	code := r.FormValue("code")
	userId := r.FormValue("state")

	t := CreateTransport(c, nil)

	user := User{}
	userKey := datastore.NewKey(c, "User", userId, 0, nil)

	err := datastore.Get(c, userKey, &user)
	if err != nil {
		c.Errorf(err.Error())
		http.Error(w, "User does not exist.", http.StatusNotFound)
		return
	}

	token, err := t.Exchange(code)
	if err != nil {
		c.Errorf(err.Error())
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	user.MovesToken = *NewMonkeyToken(token)
	user.AuthorizedWithMoves = true

	profile, err := GetUserProfile(t)
	if err != nil {
		c.Errorf(err.Error())
		http.Error(w, "Error fetching user profile from Moves.",
			http.StatusInternalServerError)
	}

	user.MovesUserId = profile.UserId

	_, err = datastore.Put(c, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
