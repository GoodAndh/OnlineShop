package exception

import (
	"net/http"
)

type NotFoundHandle func(w http.ResponseWriter, r *http.Request, v ...any) error

func NewNotFound(n NotFoundHandle) Handle {
	return func(w http.ResponseWriter, r *http.Request) {
		NotFound(w, r, map[string]interface{}{
			"error": ErrNotFoundRouter.Error(),
		})
	}
}

func NotFound(w http.ResponseWriter, r *http.Request, v ...any) error {
	return WriteNotFoundError(w, v...)
}
