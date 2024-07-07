package auth

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, SessionID string) {
	cookie := http.Cookie{
		Name:     "SessionID",
		Value:    SessionID,
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}
