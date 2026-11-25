package gastos

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const msgDatosInvalidos = "Datos inválidos, revisá los campos"

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// obtenerUsuarioID devuelve el email del usuario autenticado por la cookie de sesión.
// Lo fija el middleware security.RequerirSesion; ya no se lee ningún header del cliente.
func obtenerUsuarioID(c *gin.Context) string {
	return c.GetString("usuarioEmail")
}

// errorInterno guarda el detalle en el log del servidor y responde algo genérico.
func errorInterno(c *gin.Context, err error) {
	slog.ErrorContext(c.Request.Context(), "error en gastos", "ruta", c.FullPath(), "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrió un error interno, intentá de nuevo"})
}

func (h *Handler) CrearGasto(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)

	var in GastoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": msgDatosInvalidos})
		return
	}

	gastoNuevo := in.aGasto(0)
	id, err := h.repo.Crear(usuarioID, gastoNuevo)
	if err != nil {
		errorInterno(c, err)
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
		errorInterno(c, err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": msgDatosInvalidos})
		return
	}

	gastoEditado := in.aGasto(int64(id))
	if err := h.repo.Actualizar(usuarioID, id, gastoEditado); err != nil {
		if errors.Is(err, ErrNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		errorInterno(c, err)
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

	if err := h.repo.Borrar(usuarioID, id); err != nil {
		if errors.Is(err, ErrNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		errorInterno(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "gasto eliminado correctamente"})
}

func (h *Handler) ObtenerResumen(c *gin.Context) {
	usuarioID := obtenerUsuarioID(c)

	items, total, err := h.repo.ObtenerResumen(usuarioID)
	if err != nil {
		errorInterno(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}
