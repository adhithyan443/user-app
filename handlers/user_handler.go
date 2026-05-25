package handlers

import (
	"log/slog"
	"net/http"
	"regexp"
	"user-app/config"
	"user-app/models"
	"user-app/utils"

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

	err := config.DB.Select("id,name,email").First(&user, id).Error
	if err != nil {
		slog.Error("Failed to fetch user profile", "user_id", id, "error", err)
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
		// session.Set("form_name", name)
		// session.Set("form_email", email)
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	if len(name) < 3 {
		session.Set("message", "Name must be at least 3 characters")
		// session.Set("form_name", name)
		// session.Set("form_email", email)
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(name) {
		session.Set("message", "Name should contain only letters")
		// session.Set("form_name", name)
		// session.Set("form_email", email)
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/profile")
		return
	}

	if err := config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":  name,
			"email": email,
		}).Error; err != nil {
		slog.Error("Failed to update user profile", "user_id", id, "error", err)
		ctx.String(http.StatusInternalServerError, "Update Fail")
		return
	}

	slog.Info("User profile updated successfully", "user_id", id)

	session.Set("message", "Profile updated successfully")
	session.Set("name", name)
	session.Set("email", email)
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/profile")

}

func ShowChangePasswordPage(ctx *gin.Context) {

	session := sessions.Default(ctx)

	msg := session.Get("message")
	errmsg := session.Get("error")
	session.Delete("message")
	session.Delete("error")
	session.Save()

	ctx.HTML(http.StatusOK, "user_changepassword.html", gin.H{
		"message": msg,
		"error":   errmsg,
	})
}

func ChangePassword(ctx *gin.Context) {
	session := sessions.Default(ctx)

	id := session.Get("user_id")

	oldPassword := ctx.PostForm("oldpassword")
	newPassword := ctx.PostForm("newpassword")
	confirmPassword := ctx.PostForm("confirmpassword")

	//Pass validation
	if !utils.IsStrongPassword(newPassword) {
		session.Set("error", "Password must contain uppercase, lowercase, number, and special character")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	if newPassword != confirmPassword {
		session.Set("error", "Password do not match")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	var user models.User

	err := config.DB.Select("password").First(&user, id).Error

	if err != nil {
		slog.Error("User not found during password change", "user_id", id, "error", err)
		session.Set("error", "User not found")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		slog.Warn("Incorrect current password attempt", "user_id", id)
		session.Set("error", "Current password is incorrect")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash new password", "error", err)
		session.Set("error", "Failed to process password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	err = config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("password", string(newHashedPassword)).Error

	if err != nil {
		slog.Error("Failed to update password in database", "user_id", id, "error", err)
		session.Set("error", "Failed to update password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/password")
		return
	}

	slog.Info("User changed password successfully", "user_id", id)
	session.Set("message", "Password updated successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/password")

}
