package handlers

import (
	"log/slog"
	"net/http"
	"regexp"
	"user-app/config"
	models "user-app/internal/domain"
	"user-app/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}

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
		return
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

	//find user by email
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		slog.Warn("Login Failed - User not found", "email", email)
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		slog.Warn("Login failed - wrong password", "email", email)
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {

		slog.Error("Failed to generate JWT token",
			"user_id", user.ID,
			"error", err,
		)

		ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Failed to login",
		})

		return
	}

	//session setup
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID)
	session.Set("email", user.Email)
	session.Set("name", user.Name)
	session.Set("role", user.Role)
	session.Set("token", token)
	

	if err := session.Save(); err != nil {
		slog.Error("Failed to save session", "error", err)
	}

	slog.Info("User logged in successfully", "user_id", user.ID, "role", user.Role)
	slog.Info("JWT token generated",
		"user_id", user.ID,
		"token", token,
	)

	//Redirect based on role
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

	var count int64

	err := config.DB.Model(&models.User{}).Count(&count).Error
	if err != nil {
		slog.Error("Failed to count users", "error", err)
		ctx.String(http.StatusInternalServerError, "Error fetching user count")
		return
	}

	ctx.HTML(http.StatusOK, "admin.html", gin.H{
		"Title": "Admin Dashboard",
		"User":  "Admin User",
		"Count": count,
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
			"name":  name,
			"email": email,
		})
		return
	}

	if len(name) < 3 {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Name must be at least 3 characters",
			"name":  name,
			"email": email,
		})
		return
	}

	var nameRegex = regexp.MustCompile(`^[a-zA-Z ]+$`)

	if !nameRegex.MatchString(name) {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Name should contain only letters",
			"name":  name,
			"email": email,
		})
		return
	}

	if !utils.IsStrongPassword(password) {
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Password must contain uppercase, lowercase, number, and special character",
			"name":  name,
			"email": email,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Failed to hash password")
		ctx.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		slog.Warn("Signup failed - duplicate email?", "email", email, "error", err)
		ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Email already exists or invalid data",
			"name":  name,
			"email": email,
		})
		return
	}

	slog.Info("New user registered", "user_id", user.ID, "email", email)

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

		var user models.User

		err := config.DB.Where("email = ?", email).First(&user).Error
		if err != nil {
			slog.Warn("Forgot password - user not found", "email", email)
			session.Set("message", "User not exsists")
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/forgotpassword")
			return
		}

		session.Set("reset_id", user.ID)
		session.Save()
		ctx.HTML(http.StatusOK, "forgotpassword.html", gin.H{
			"account": false,
		})

	} else {

		id := session.Get("reset_id")
		newpass := ctx.PostForm("newpassword")
		confirmpass := ctx.PostForm("confirmpassword")

		if !utils.IsStrongPassword(newpass) {
			session.Set("message", "Password must contain uppercase, lowercase, number, and special character")
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/forgotpassword")
			return
		}

		if newpass != confirmpass {
			session.Set("message", "Password do not match")
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/forgotpassword")
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Failed to hash new password")
			ctx.String(http.StatusInternalServerError, "Error hashing password")
			return
		}

		err = config.DB.Model(&models.User{}).
			Where("id = ?", id).
			Update("password", string(hashed_password)).Error

		if err != nil {
			slog.Error("Failed to update password", "error", err)
			ctx.String(http.StatusInternalServerError, "Database error")
			return
		}

		session.Delete("reset_id")
		session.Set("message", "Password updated successfully")
		session.Save()

		slog.Info("Password reset successful", "user_id", id)
		ctx.Redirect(http.StatusSeeOther, "/login")

	}
}

func JWTProfile(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")
	email, _ := ctx.Get("email")
	role, _ := ctx.Get("role")

	slog.Info("JWT protected route accessed",
		"user_id", userID,
		"email", email,
		"role", role,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "JWT authentication successful",
		"user_id": userID,
		"email": email,
		"role": role,
	})
}