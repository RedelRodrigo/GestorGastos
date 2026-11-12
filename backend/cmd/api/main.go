package main

import (
	"log"
	"os"

	"GestorGastos/internal/db"
	"GestorGastos/internal/gastos"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB(getenv("DB_PATH", "./gastos.db"))
	defer database.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api funcionando"})
	})

	gastos.RouterGasto(router, database)

	port := getenv("PORT", "8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
