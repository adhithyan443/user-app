package routes

import (
	"user-app/handlers"
	"user-app/internal/handler"
	"user-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler,userHandler *handler.UserHandler,adminHandler *handler.AdminHandler,  pageHandler *handler.PageHandler) {

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(303, "/login")
	})

	//Login Routes
	r.GET("/login", pageHandler.ShowLoginPage)
	r.POST("/login", authHandler.HandleLogin)
	//Signup
	r.GET("/signup", pageHandler.ShowSignupPage)
	r.POST("/signup", authHandler.HandleSignup)

	r.GET("/forgotpassword", pageHandler.ShowForgotPasswordPage)
	r.POST("/forgotpassword", authHandler.HandleForgotPassword)

	//middleware
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(), middleware.NoCache())
	{
		protected.GET("/home", handlers.ShowHomePage)
		protected.GET("/logout", authHandler.Logout)

		//User route
		protected.GET("/profile", 	userHandler.ShowProfilePage)
		protected.POST("/profile/update", userHandler.UpdateUserProfile)

		protected.GET("/password", pageHandler.ShowChangePasswordPage)
		protected.POST("/password", userHandler.ChangePassword)

		//admin route
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("", handlers.ShowAdminPage)                                 //dashboard

			admin.GET("/users",  adminHandler.GetAllUsers)                              //Read all users
			admin.GET("/users/edit/:id", adminHandler.EditUserPage)                   //Edit user
			admin.POST("/users/update/:id", adminHandler.UpdateUserPage)              //Update
			admin.GET("/users/delete/:id", adminHandler.DeleteUser)                   //Delete user

			admin.GET("/users/updatepassword/:id", adminHandler.ShowUserPasswordPage) //update password
			admin.POST("/users/updatepassword/:id", adminHandler.EditUserPasswordPage)

			admin.GET("/newuser", adminHandler.NewUserPage)
			admin.POST("/newuser", adminHandler.AddNewUser)

		}

	}

	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/profile", handlers.JWTProfile)
	}
}
