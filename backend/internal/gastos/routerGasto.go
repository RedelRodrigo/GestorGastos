package gastos

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RouterGasto(router *gin.Engine, database *sql.DB) {

	router.POST("/gastos", crearGastosHandler(database))

	router.GET("/gastos", obtenerGastosHandler(database))

	router.PUT("/gastos/:id", editarGastosHandler(database))

	router.DELETE("/gastos/:id", borrarGastoHandler(database))
}
