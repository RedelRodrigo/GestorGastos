package gastos

import (
	"database/sql"

	"GestorGastos/internal/security"

	"github.com/gin-gonic/gin"
)

// RouterGasto registra las rutas de /api/gastos. Todas exigen sesión válida:
// sin cookie firmada, el middleware responde 401 y el handler nunca se ejecuta.
func RouterGasto(rg *gin.RouterGroup, db *sql.DB) {
	repo := NewSQLiteRepository(db)
	h := NewHandler(repo)

	gastosGroup := rg.Group("/gastos", security.RequerirSesion())
	{
		gastosGroup.POST("", h.CrearGasto)
		gastosGroup.GET("", h.ListarGastos)
		gastosGroup.GET("/resumen", h.ObtenerResumen)
		gastosGroup.PUT("/:id", h.EditarGasto)
		gastosGroup.DELETE("/:id", h.BorrarGasto)
	}
}
