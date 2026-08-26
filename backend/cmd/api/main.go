package main

import (
	"GestorGastos/internal/db"
	"GestorGastos/internal/gastos"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database := db.InitDB()
	defer database.Close()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api funcionando",
		})
	})
	gastos.RouterGasto(router, database)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
