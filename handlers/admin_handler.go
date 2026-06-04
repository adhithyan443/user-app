package handlers

import (

	"log/slog"
	"net/http"
	"strconv"
	"user-app/config"
	models "user-app/internal/domain"
	"user-app/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	// "golang.org/x/text/message"
)

func GetAllUser(ctx *gin.Context) {

	session := sessions.Default(ctx)
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	search := ctx.Query("search")

	var user []models.User
	var err error

	if search != "" {
		err = config.DB.Where("name ILike ?", "%"+search+"%").
			Or("email ILIKE ?", "%"+search+"%").
			Find(&user).Error

	} else {
		err = config.DB.Find(&user).Error
	}

	if err != nil {
		slog.Error("Failed to fetch users", "error", err, "search", search)
		ctx.String(http.StatusInternalServerError, "Error feching users")
		return
	}

	ctx.HTML(http.StatusOK, "admin_users.html", gin.H{
		"users":   user,
		"message": msg,
	})
}



func ShowUserPasswordPage(ctx *gin.Context) {

	session := sessions.Default(ctx)
	id := ctx.Param("id")
	msg := session.Get("message")
	session.Delete("message")
	session.Save()
	ctx.HTML(http.StatusOK, "admin_changepassword.html", gin.H{
		"message": msg,
		"id":      id,
	})
}

func EditUserPasswordPage(ctx *gin.Context) {

	session := sessions.Default(ctx)

	newpassword := ctx.PostForm("newpassword")
	confirmpassword := ctx.PostForm("confirmpassword")
	id := ctx.Param("id")

	if len(newpassword) < 6 {
		session.Set("message", "Password must be at least 6 characters")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	if !utils.IsStrongPassword(newpassword) {
		session.Set("message", "Password must contain uppercase, lowercase, number, and special character")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	if newpassword != confirmpassword {
		session.Set("message", "Passwords do not match")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		session.Set("message", "Failed to process password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	err = config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("password", string(hashedPassword)).Error

	if err != nil {
		slog.Error("Failed to update user password", "user_id", id, "error", err)
		session.Set("message", "Failed to update password")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+id)
		return
	}

	slog.Info("Admin updated user password successfully", "user_id", id)
	session.Set("message", "Password updated successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
}

func DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID format for delete", "id_param", idParam)
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = config.DB.Unscoped().Delete(&models.User{}, id).Error

	if err != nil {
		slog.Error("Failed to delete user", "user_id", id, "error", err)
		ctx.String(http.StatusInternalServerError, "Delete failed")
		return
	}

	slog.Info("User deleted successfully", "user_id", id)

	session := sessions.Default(ctx)
	session.Set("message", "User deleted successfully")
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/users")
}


