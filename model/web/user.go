package web

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterPayload struct {
	Username    string `json:"username" validate:"required,min=8,max=30"`
	Password    string `json:"password" validate:"required,min=8"`
	CekPassword string `json:"c_password" validate:"required,eqfield=Password"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name" validate:"required"`
}
