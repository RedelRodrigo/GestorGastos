package gastos

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Cambiamos el parámetro para recibir el group /api directamente
func RouterGasto(rg *gin.RouterGroup, db *sql.DB) {
	// Nota: Ya no necesitamos configurar CORS aquí si ya lo maneja el router principal,
	// o puedes dejarlo si prefieres.

	repo := NewSQLiteRepository(db)
	h := NewHandler(repo)

	// Como rg ya viene con "/api", aquí solo definimos "/gastos" (quedando /api/gastos)
	gastosGroup := rg.Group("/gastos")
	{
		gastosGroup.POST("", h.CrearGasto)
		gastosGroup.GET("", h.ListarGastos)
		gastosGroup.GET("/resumen", h.ObtenerResumen)
		gastosGroup.PUT("/:id", h.EditarGasto)
		gastosGroup.DELETE("/:id", h.BorrarGasto)
	}
}
