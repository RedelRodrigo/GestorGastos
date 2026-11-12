package gastos

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterGasto(r *gin.Engine, db *sql.DB) {
	// Configuración de CORS para el frontend en Vite (puerto 5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "X-User-ID"}, // <- Añadido X-User-ID
	}))

	// Inyección de dependencias
	repo := NewSQLiteRepository(db)
	h := NewHandler(repo)

	// Agrupación de rutas
	gastosGroup := r.Group("/gastos")
	{
		gastosGroup.POST("", h.CrearGasto)
		gastosGroup.GET("", h.ListarGastos)
		gastosGroup.GET("/resumen", h.ObtenerResumen)
		gastosGroup.PUT("/:id", h.EditarGasto)
		gastosGroup.DELETE("/:id", h.BorrarGasto)
	}
}
