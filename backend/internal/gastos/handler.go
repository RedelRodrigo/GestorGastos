package gastos

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// obtenerUsuarioID extrae el identificador del usuario desde los headers HTTP.
func obtenerUsuarioID(c *gin.Context) string {
	usuarioID := c.GetHeader("X-User-ID")
	if usuarioID == "" {
		return "usuario_invitado" // Fallback seguro
	}
	return usuarioID
}

func (h *Handler) CrearGasto(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)
	var in GastoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gastoNuevo := in.aGasto(0)
	id, err := h.repo.Crear(usuarioID, gastoNuevo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gastoNuevo.ID = id
	c.JSON(http.StatusCreated, gin.H{"data": gastoNuevo})
}

func (h *Handler) ListarGastos(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)
	filtro := Filtro{
		Categoria: c.Query("categoria"),
		Desde:     c.Query("desde"),
		Hasta:     c.Query("hasta"),
	}

	gastos, err := h.repo.Listar(usuarioID, filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gastos})
}

func (h *Handler) EditarGasto(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var in GastoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gastoEditado := in.aGasto(int64(id))
	err = h.repo.Actualizar(usuarioID, id, gastoEditado)
	if err != nil {
		if errors.Is(err, ErrNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gastoEditado})
}

func (h *Handler) BorrarGasto(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	err = h.repo.Borrar(usuarioID, id)
	if err != nil {
		if errors.Is(err, ErrNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "gasto eliminado correctamente"})
}

func (h *Handler) ObtenerResumen(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)
	items, total, err := h.repo.ObtenerResumen(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}
