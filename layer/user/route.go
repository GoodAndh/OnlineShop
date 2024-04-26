package user

import (
	"ddd2/exception"
	"ddd2/layer/auth"
	"ddd2/model/web"
	"ddd2/utils"
	"ddd2/views"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Validate *validator.Validate
	Service  Service
}

func NewRoute(validate *validator.Validate, service Service) *Route {
	return &Route{
		Validate: validate,
		Service:  service,
	}
}

func (h *Route) RegisterRoute(router *httprouter.Router) {
	router.GET("/login/", h.handleLogin)
	router.POST("/login/", h.handleLogin)

	router.GET("/register/", h.handleRegister)
	router.POST("/register/", h.handleRegister)
}

func (h *Route) handleLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		views.TemplateExecuted(w, nil, "login.html")
		return
	case http.MethodPost:
		r.ParseForm()
		errorList := make(map[any]any)

		user := &web.UserLoginPayload{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		u, err := h.Service.GetByUsername(r.Context(), user.Username)
		if err != nil {
			errorList["error"] = exception.ErrIncorrectUser
			views.TemplateExecuted(w, errorList, "login.html")
			return
		}

		if !auth.ComparePassword(u.Password, []byte(user.Password)) {
			errorList["error"] = exception.ErrIncorrectUser
			views.TemplateExecuted(w, errorList, "login.html")
			return
		}

		errorList["data"] = u
		var dataList int

		if data, ok := errorList["data"].(*web.User); ok {
			dataList=data.Id
		}

		ses, _ := auth.Session.Get(r, "lg-ses")
		err = auth.SaveSession(w, r, ses, dataList)
		if err != nil {
			errorList["err"]=err
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		time.Sleep(2*time.Second)
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
}

func (h *Route) handleRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		views.TemplateExecuted(w, nil, "register.html")
		return
	case http.MethodPost:
		r.ParseForm()

		payload := &web.UserRegisterPayload{
			Username:    r.Form.Get("username"),
			Password:    r.Form.Get("password"),
			CekPassword: r.Form.Get("cpassword"),
			Email:       r.Form.Get("email"),
			Name:        r.Form.Get("name"),
		}

		errorList, err := utils.ValidationField(h.Validate, payload)
		if err != nil {
			views.TemplateExecuted(w, errorList, "register.html")
			return
		}

		err = h.Service.CreateUser(r.Context(), *payload)
		if err != nil {
			if errors.Is(err, exception.ErrEmailTaken) {
				errorList["error"] = exception.ErrEmailTaken
				views.TemplateExecuted(w, errorList, "register.html")
				return
			}
			if errors.Is(err, exception.ErrUsernameTaken) {
				errorList["error"] = exception.ErrEmailTaken
				views.TemplateExecuted(w, errorList, "register.html")
				return
			}

			errorList["error"] = err
			views.TemplateExecuted(w, errorList, "register.html")
			return
		}

	}
}
