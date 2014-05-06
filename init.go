package amelia

import (
	"net/http"

	"appengine"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.ERror.
		c := appengine.NewContext(r)
		c.Errorf("%v", e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func init() {
	http.Handle("/", appHandler(root))
	http.Handle("/login", appHandler(login))
	http.Handle("/logout", appHandler(logout))
	http.Handle("/addphone", appHandler(addPhone))
	http.Handle("/delphone", appHandler(delPhone))
	http.Handle("/authorize", appHandler(authorize))
	http.Handle("/oauth2callback", appHandler(oauthCallback))
	http.Handle("/revoke", appHandler(revoke))
	http.Handle("/notification", appHandler(handleNotification))
}
