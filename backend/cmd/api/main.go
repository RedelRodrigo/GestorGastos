package main

import (
	"log"
	"os"

	"GestorGastos/internal/auth"
	"GestorGastos/internal/db"
	"GestorGastos/internal/gastos"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de entorno del sistema.")
	}

	database := db.InitDB()
	defer database.Close()

	router := gin.Default()

	// ==========================================
	// 1. RUTAS DE LA API (Deben ir PRIMERO)
	// ==========================================
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "api funcionando"})
		})

		authRepo := auth.NewSQLiteAuthRepository(database)
		authHandler := auth.NewHandler(authRepo)
		api.POST("/auth/registro", authHandler.Registrar)
		api.POST("/auth/login", authHandler.Login)

		// 🟢 Registramos las rutas de gastos dentro del grupo /api
		gastos.RouterGasto(api, database)
	}

	// ==========================================
	// 2. CONFIGURACIÓN DEL FRONTEND ESTÁTICO (SPA)
	// ==========================================
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}

	// Servir assets de Vite
	router.Static("/assets", staticDir+"/assets")

	// Servir index.html en la raíz
	router.StaticFile("/", staticDir+"/index.html")

	// SPA Fallback: Cualquier ruta que no sea /api/... ni un asset físico
	// devolverá el index.html para que React Router controle las vistas.
	router.NoRoute(func(c *gin.Context) {
		c.File(staticDir + "/index.html")
	})

	// ==========================================
	// 3. INICIO DEL SERVIDOR
	// ==========================================
	port := getenv("PORT", "8080")
	log.Printf("Servidor corriendo en puerto %s", port)
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
