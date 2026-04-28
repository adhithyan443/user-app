package handlers

import (
	"net/http"
	"regexp"
	"user-app/config"
	"user-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowLoginPage(ctx *gin.Context) {

	session := sessions.Default(ctx)
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	if session.Get("user_id") != nil {
		if session.Get("role") == "admin" {
			ctx.Redirect(http.StatusSeeOther, "/admin")
		} else {
			ctx.Redirect(http.StatusSeeOther, "/home")
		}

	}

	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"message": msg,
	})
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
		"Count": userCount,
	})
}

func ShowSignupPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

func HandleSignup(ctx *gin.Context) {

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	//Validation

	if name == "" || email == "" || password == "" {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "All fields are required",
		})
		return
	}

	if len(name) < 3 {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Name must be at least 3 characters",
		})
		return
	}

	var nameRegex = regexp.MustCompile(`^[a-zA-Z ]+$`)

	if !nameRegex.MatchString(name) {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Name should contain only letters",
		})
		return
	}

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

	session := sessions.Default(ctx)
	session.Set("message", "Profile created successfully")
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/login")
}

func HandleLogout(ctx *gin.Context) {

	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusFound, "/login")
}

func ShowForgotPasswordPage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	ctx.HTML(http.StatusOK, "forgotpassword.html", gin.H{
		"account": true,
		"message": msg,
	})
}

func HandleForgotPassword(ctx *gin.Context) {
	session := sessions.Default(ctx)
	email := ctx.PostForm("email")
	// id:=ctx.Param("id")

	if email != "" {
		var user_id int
		err := config.DB.QueryRow("SELECT id FROM users WHERE email =$1", email).Scan(&user_id)

		if err != nil {

			session.Set("message", "User not exsists")
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/forgotpassword")
			return
		}

		session.Set("reset_id", user_id)
		session.Save()
		ctx.HTML(http.StatusOK, "forgotpassword.html", gin.H{
			"account": false,
		})

	} else {
		id := session.Get("reset_id")
		newpass := ctx.PostForm("newpassword")
		confirmpass := ctx.PostForm("confirmpassword")

		if newpass != confirmpass {
			session.Set("message", "Password do not match")
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/forgotpassword")
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Error hashing password")
			return
		}

		_, err = config.DB.Exec("UPDATE users SET hashed_password = $1 WHERE id=$2", hashed_password, id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Database error")
			return
		}

		session.Delete("reset_id")
		session.Set("message", "Password updated successfully")
		session.Save()

		ctx.Redirect(http.StatusSeeOther, "/login")

	}
}
