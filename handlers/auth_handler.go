package handlers

import (
	"net/http"
	"user-app/config"
	"user-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

// Login Handler
func HandleLogin(ctx *gin.Context) {

	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	var user models.User

	//1. Get user from db
	query := `SELECT id,name,hashed_password, role FROM users WHERE email=$1`

	err := config.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.HashedPassword, &user.Role)

	//check user exists
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//3. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//session setup
	session := sessions.Default(ctx)

	session.Set("user_id", user.ID)
	session.Set("email", user.Email)
	session.Set("name", user.Name)
	session.Set("role", user.Role)
	session.Save()

	//4.Redirect based on role
	if user.Role == "admin" {
		ctx.Redirect(http.StatusSeeOther, "/admin")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/home")
	}
}

func ShowHomePage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	name := session.Get("name")

	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"name": name,
	})
}

func ShowAdminPage(ctx *gin.Context) {

	query := `SELECT COUNT(*) FROM users`

	var userCount int

	err := config.DB.QueryRow(query).Scan(&userCount)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error fetching user count")
		return
	}

	ctx.HTML(http.StatusOK, "admin.html", gin.H{
		"Title": "Admin Dashboard",
		"User":  "Admin User",
		"Count":userCount,
	})
}

func ShowSignupPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

func HandleSignup(ctx *gin.Context) {

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	query := `INSERT INTO users(name,email,hashed_password) VALUES ($1,$2,$3)`
	_, err = config.DB.Exec(query, name, email, hashedPassword)

	if err != nil {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Email already exists or invalid data",
		})
	}

	ctx.Redirect(http.StatusSeeOther, "/login")
}

func HandleLogout(ctx *gin.Context) {

	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusFound, "/login")
}
