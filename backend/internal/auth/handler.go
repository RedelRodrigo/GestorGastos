package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"GestorGastos/internal/security"

	"github.com/gin-gonic/gin"
)

const msgCaptcha = "No pudimos verificar que sos una persona. Intentá de nuevo."

type Handler struct {
	repo      Repository
	porIP     *security.LimitadorLogin
	porCorreo *security.LimitadorLogin
}

func NewHandler(repo Repository, porIP, porCorreo *security.LimitadorLogin) *Handler {
	return &Handler{repo: repo, porIP: porIP, porCorreo: porCorreo}
}

func (h *Handler) Registrar(c *gin.Context) {
	ctx := c.Request.Context()

	var in RegistroInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos. La contraseña debe tener al menos 12 caracteres."})
		return
	}
	if !security.VerificarTurnstile(ctx, in.Captcha, c.ClientIP()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": msgCaptcha})
		return
	}

	usuario, err := h.repo.Registrar(in)
	if err != nil {
		if errors.Is(err, ErrEmailDuplicado) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		slog.ErrorContext(ctx, "registro", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al registrar"})
		return
	}

	if err := security.EmitirSesion(c.Writer, usuario.Email, usuario.Username); err != nil {
		slog.ErrorContext(ctx, "emitir sesión", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al registrar"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": usuario})
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var in LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingresá un correo y contraseña válidos"})
		return
	}

	ip := "ip:" + c.ClientIP()
	correo := "correo:" + strings.ToLower(strings.TrimSpace(in.Email))

	// Se revisa el límite antes de llamar a Turnstile, para no gastar verificaciones.
	if !h.porIP.Permitir(ip) || !h.porCorreo.Permitir(correo) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Demasiados intentos. Esperá unos minutos e intentá de nuevo."})
		return
	}
	if !security.VerificarTurnstile(ctx, in.Captcha, c.ClientIP()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": msgCaptcha})
		return
	}

	usuario, err := h.repo.Login(in)
	if err != nil {
		if errors.Is(err, ErrCredencialesInv) {
			h.porIP.Fallo(ip)
			h.porCorreo.Fallo(correo)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		slog.ErrorContext(ctx, "login", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al iniciar sesión"})
		return
	}

	h.porCorreo.Exito(correo)
	if err := security.EmitirSesion(c.Writer, usuario.Email, usuario.Username); err != nil {
		slog.ErrorContext(ctx, "emitir sesión", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al iniciar sesión"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": usuario})
}

func (h *Handler) Sesion(c *gin.Context) {
	id := c.MustGet("usuario").(*security.Identidad)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"email": id.Email, "username": id.Username}})
}

func (h *Handler) Logout(c *gin.Context) {
	security.CerrarSesion(c.Writer)
	c.JSON(http.StatusOK, gin.H{"mensaje": "sesión cerrada"})
}