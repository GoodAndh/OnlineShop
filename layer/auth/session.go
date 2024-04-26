package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	randomString = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890-=/!@#$%^&*()_+"
)

// func generateRandomString(size int, charset string) string {
// 	b := make([]byte, size)
// 	for i := range b {
// 		b[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(b)
// }
// func generateString(size int) string {
// 	return generateRandomString(size, randomString)
// }

var (
	Session = sessions.NewCookieStore([]byte(randomString))
)

func SaveSession(w http.ResponseWriter, r *http.Request, ses *sessions.Session, data any) error {
	ses.Values["auten"] = true
	ses.Values["data"] = data
	ses.Options.MaxAge = 3600
	return ses.Save(r, w)
}
