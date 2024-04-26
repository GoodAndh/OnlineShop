package exception

import (
	"ddd2/model/web"
	"ddd2/views"
	"encoding/json"
	"log"
	"net/http"
)

func ParseJson(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJson(w http.ResponseWriter, code int, status, message string, data ...any) error {
	web := &web.WebResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
	return json.NewEncoder(w).Encode(web)
}


func WriteInternalServerError(w http.ResponseWriter, mp map[any]any, filename ...string) error {

	if noSes := mp["error"].(error).Error(); noSes == ErrNoSesFound.Error() {
		mp["errormsg"] = true
		log.Println(mp)
		return views.TemplateExecuted(w, mp, filename...)

	} else {
		return views.TemplateExecuted(w, mp, filename...)
	}

}

func WriteNotFoundError(w http.ResponseWriter, data ...any) error {
	return WriteJson(w, http.StatusNotFound, "not found error", "not found url/path", data...)
}
