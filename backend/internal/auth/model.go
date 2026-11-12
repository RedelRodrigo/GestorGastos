package auth

import "errors"

var (
	ErrEmailDuplicado  = errors.New("el correo electrónico ya está registrado")
	ErrCredencialesInv = errors.New("correo o contraseña incorrectos")
)

type Usuario struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RegistroInput struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=40"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
