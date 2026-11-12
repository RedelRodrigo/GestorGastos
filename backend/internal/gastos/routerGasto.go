package gastos

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterGasto(router *gin.Engine, database *sql.DB) {

	// Dentro de tu inicialización de rutas:
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	router.POST("/gastos", crearGastosHandler(database))

	router.GET("/gastos", obtenerGastosHandler(database))

	router.GET("/gastos/:id", obtenerGastoHandlerId(database))

	router.PUT("/gastos/:id", editarGastosHandler(database))

	router.DELETE("/gastos/:id", borrarGastoHandler(database))

	router.GET("/gastos/resumen", resumenHandler(database))
}
