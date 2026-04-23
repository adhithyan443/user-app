package main

import (
	"fmt"
	"html/template"
	"net/http"

	"user-app/config"
	"user-app/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase() // -> Connect to database
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

	fmt.Println("🔴Server is running on http://localhost:8080")
	fmt.Println("Session middleware enabled")

	r.Run(":8080")
}
