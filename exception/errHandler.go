package exception

import "net/http"

type BaseError interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Handle func(http.ResponseWriter, *http.Request)

func (n Handle)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	n(w,r)
}
