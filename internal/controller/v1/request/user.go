package request

type Login struct {
	Password string `json:"password" validate:"required"`
} // @name v1.Login
