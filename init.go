package gamehack

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/addphone", addPhone)
	http.HandleFunc("/delphone", delPhone)
	http.HandleFunc("/authorize", authorize)
	http.HandleFunc("/oauth2callback", oauthCallback)
	http.HandleFunc("/revoke", revoke)
	http.HandleFunc("/notification", handleNotification)
}
