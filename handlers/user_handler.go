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

func ShowProfilePage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	id := session.Get("user_id")
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	var user models.User

	query := `SELECT name,email from users WHERE id = $1`

	err := config.DB.QueryRow(query, id).Scan(&user.Name, &user.Email)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "User not found")
		return
	}

	ctx.HTML(http.StatusOK, "user_profile.html", gin.H{
		"user":    user,
		"message": msg,
	})
}

func UpdateUserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	id := session.Get("user_id")

	name := ctx.PostForm("name")
	email := ctx.PostForm("email")

	if name == "" || email == "" {
		session.Set("message", "All fields are required")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	if len(name) < 3 {
		session.Set("message", "Name must be at least 3 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(name) {
		session.Set("message", "Name should contain only letters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	
	_, err := config.DB.Exec(
		"UPDATE users SET name=$1,email=$2 WHERE id=$3", name, email, id,
	)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Update Fail")
		return
	}

	session.Set("message", "Profile updated successfully")
	session.Set("name", name)
	session.Set("email", email)
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/profile")

}

func ShowChangePasswordPage(ctx *gin.Context) {

	session:=sessions.Default(ctx)

	msg:=session.Get("message")
	errmsg:=session.Get("error")
	session.Delete("message")
	session.Delete("error")
	session.Save()

	ctx.HTML(http.StatusOK, "user_changepassword.html", gin.H{
		"message":msg,
		"error":errmsg,
	})
}

func ChangePassword(ctx *gin.Context) {
	session := sessions.Default(ctx)

	id := session.Get("user_id")

	oldPassword := ctx.PostForm("oldpassword")
	newPassword := ctx.PostForm("newpassword")
	confirmPassword := ctx.PostForm("confirmpassword")

	if newPassword != confirmPassword {
		session.Set("error", "Password do not match")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	var hashedpassword string

	err := config.DB.QueryRow(
		"SELECT hashed_password FROM users WHERE id=$1", id,
	).Scan(&hashedpassword)

	if err != nil {
		session.Set("error", "User not found")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(oldPassword))
	if err != nil {
		session.Set("error", "Current password is incorrect")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		session.Set("error", "Failed to process password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	_, err = config.DB.Exec(
		"UPDATE users SET hashed_password=$1 WHERE id=$2",
		newHashedPassword,
		id,
	)

	if err != nil {
		session.Set("error", "Failed to update password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	session.Set("message", "Password updated successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/password")

}
