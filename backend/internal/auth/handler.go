package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Registrar(c *gin.Context) {
	var in RegistroInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos (contraseña mínima de 6 caracteres)"})
		return
	}

	usuario, err := h.repo.Registrar(in)
	if err != nil {
		if errors.Is(err, ErrEmailDuplicado) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al registrar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": usuario})
}

func (h *Handler) Login(c *gin.Context) {
	var in LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingresá un correo y contraseña válidos"})
		return
	}

	usuario, err := h.repo.Login(in)
	if err != nil {
		if errors.Is(err, ErrCredencialesInv) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al iniciar sesión"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": usuario})
}
