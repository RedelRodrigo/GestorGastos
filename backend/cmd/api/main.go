package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"GestorGastos/internal/auth"
	"GestorGastos/internal/db"
	"GestorGastos/internal/gastos"
	"GestorGastos/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("no se encontró .env, usando variables del sistema")
	}
	if err := security.InicializarSecreto(); err != nil {
		slog.Error("configuración inválida", "error", err)
		os.Exit(1)
	}

	database := db.InitDB()
	defer database.Close()

	router := gin.Default()
	// IP del proxy (Caddy) que está delante de la app. Solo esas IPs pueden fijar X-Forwarded-For.
	if err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		slog.Error("proxies confiables inválidos", "error", err)
		os.Exit(1)
	}
	router.Use(security.Headers())

	porIP := security.NewLimitadorLogin(20, 15*time.Minute, 15*time.Minute)
	porCorreo := security.NewLimitadorLogin(5, 15*time.Minute, 30*time.Minute)

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "api funcionando"})
		})

		authRepo := auth.NewSQLiteAuthRepository(database)
		authHandler := auth.NewHandler(authRepo, porIP, porCorreo)
		api.POST("/auth/registro", authHandler.Registrar)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/sesion", security.RequerirSesion(), authHandler.Sesion)

		gastos.RouterGasto(api, database)
	}

	staticDir := getenv("STATIC_DIR", "./static")
	router.Static("/assets", staticDir+"/assets")
	router.StaticFile("/", staticDir+"/index.html")

	// Las rutas /api inexistentes devuelven JSON; el resto cae en la SPA.
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
			return
		}
		c.File(staticDir + "/index.html")
	})

	port := getenv("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		var err error
		if cert, key := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"); cert != "" && key != "" {
			slog.Info("servidor HTTPS escuchando", "puerto", port)
			err = srv.ListenAndServeTLS(cert, key)
		} else {
			slog.Info("servidor HTTP escuchando (usa un proxy con HTTPS en producción)", "puerto", port)
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("servidor detenido", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("apagando servidor")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("apagado forzado", "error", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
