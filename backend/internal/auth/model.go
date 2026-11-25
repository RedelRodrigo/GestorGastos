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

// RegistroInput: contraseña de 8 a 72 caracteres (bcrypt no usa más de 72 bytes).
type RegistroInput struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Username string `json:"username" binding:"required,min=2,max=40"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Captcha  string `json:"captcha" binding:"required"`
}

// LoginInput no exige 12 caracteres, así las cuentas anteriores siguen pudiendo entrar.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,max=72"`
	Captcha  string `json:"captcha" binding:"required"`
}
