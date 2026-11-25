package security

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	nombreCookie   = "sesion"
	duracionSesion = 8 * time.Hour
	emisor         = "gestor-gastos"
)

var secreto []byte

type Identidad struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// InicializarSecreto carga JWT_SECRET y falla si es débil.
// Generarlo con: openssl rand -base64 48
func InicializarSecreto() error {
	s := os.Getenv("JWT_SECRET")
	if len(s) < 32 {
		return errors.New("JWT_SECRET debe tener al menos 32 caracteres")
	}
	secreto = []byte(s)
	return nil
}

// EmitirSesion guarda el token en una cookie HttpOnly, Secure y SameSite=Strict.
func EmitirSesion(w http.ResponseWriter, email, username string) error {
	ahora := time.Now()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Identidad{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    emisor,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(duracionSesion)),
		},
	}).SignedString(secreto)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     nombreCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(duracionSesion.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func CerrarSesion(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     nombreCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// RequerirSesion valida la cookie. El email del contexto es la única identidad.
func RequerirSesion() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(nombreCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tenés que iniciar sesión"})
			return
		}
		id := &Identidad{}
		_, err = jwt.ParseWithClaims(cookie.Value, id,
			func(*jwt.Token) (any, error) { return secreto, nil },
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithExpirationRequired(),
			jwt.WithIssuer(emisor),
			jwt.WithIssuedAt(),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tu sesión venció, volvé a iniciar sesión"})
			return
		}
		c.Set("usuario", id)
		c.Set("usuarioEmail", id.Email)
		c.Next()
	}
}
