package gamehack

import (
	"time"

	"code.google.com/p/goauth2/oauth"
)

// MonkeyToken is a monkeypatched oauth token that can directly be stored in
// Google Datastore. It doesn't contain a map of extra goodies.
type MonkeyToken struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time // If zero the token has no (known) expiry time.
}

func (t *MonkeyToken) Expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Before(time.Now())
}

func (t *MonkeyToken) NormToken() *oauth.Token {
	return &oauth.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
	}
}

func NewMonkeyToken(token *oauth.Token) *MonkeyToken {
	return &MonkeyToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
}
