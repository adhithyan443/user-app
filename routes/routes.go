package routes

import (
	"user-app/handlers"
	"user-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	//Login Routes
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.HandleLogin)
	//Signup
	r.GET("/signup", handlers.ShowSignupPage)
	r.POST("/signup", handlers.HandleSignup)

	r.GET("/forgotpassword",handlers.ShowForgotPasswordPage)
	r.POST("/forgotpassword",handlers.HandleForgotPassword)

	//middleware
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/home", handlers.ShowHomePage)
		protected.GET("/logout", handlers.HandleLogout)


		//User route
		protected.GET("/profile",handlers.ShowProfilePage)
		protected.POST("/profile/update",handlers.UpdateUserProfile)

		protected.GET("/password",handlers.ShowChangePasswordPage)
		protected.POST("/password",handlers.ChangePassword)


		//admin route
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("", handlers.ShowAdminPage)    //dashboard
			admin.GET("/users", handlers.GetAllUser) //Read all users
			admin.GET("/users/edit/:id",handlers.EditUserPage)//Edit user
			admin.POST("/users/update/:id",handlers.UpdateUserPage)//Update
			admin.GET("/users/delete/:id",handlers.DeleteUser)//Delete user

			admin.GET("/newuser",handlers.NewUserPage)
			admin.POST("/newuser",handlers.AddNewUser)

		}

	}
}
