package main

import (
	
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"user-app/config"
	"user-app/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	
	//Setup logger
	setupLogger()

	//Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	slog.Info("Starting User Application...")

	// Connect to database + migration
	db := config.ConnectDatabase() 
	config.AutoMigrate(db)

	r := gin.Default()

	//session setup
	stores := cookie.NewStore([]byte("super-secret-key-1234"))
	stores.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, //24hr
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("mysession", stores))

	slog.Info("Session middleware enabled")

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	//Load HTML templates
	r.LoadHTMLGlob("templates/**/*.html")
	// Serve static files (CSS, JS, images)
	r.Static("/static", "./templates/static")

	//Setup all routes
	routes.SetupRoutes(r)

	slog.Info("Server is running on http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}

func setupLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))
}
