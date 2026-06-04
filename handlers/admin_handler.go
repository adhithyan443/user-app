package handlers

import (

	"log/slog"
	"net/http"
	"regexp"
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

func NewUserPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// msg := session.Get("message")
	data := gin.H{
		"message":    session.Get("message"),
		"form_name":  session.Get("form_name"),
		"form_email": session.Get("form_email"),
		"form_role":  session.Get("form_role"),
	}

	session.Delete("form_name")
	session.Delete("form_email")
	session.Delete("form_role")
	session.Delete("message")
	session.Save()

	ctx.HTML(http.StatusOK, "admin_add_user.html", data)
	// gin.H{
	// 	// "message": msg,
	// 	"data":data,
	// })
}

func AddNewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	role := ctx.PostForm("role")
	password := ctx.PostForm("password")

	formData := map[string]interface{}{
		"form_name":  name,
		"form_email": email,
		"form_role":  role,
	}

	// Required fields
	if name == "" || email == "" || role == "" || password == "" {
		session.Set("message", "All fields are required")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Name validation (letters + space, min 3)
	if len(name) < 3 {
		session.Set("message", "Name must be at least 3 characters")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(name) {
		session.Set("message", "Name should contain only letters")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	if !utils.IsStrongPassword(password) {
		session.Set("message", "Password must contain uppercase, lowercase, number, and special character")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Password validation
	if len(password) < 6 {
		session.Set("message", "Password must be at least 6 characters")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	// Role validation
	if role != "admin" && role != "user" {
		session.Set("message", "Invalid role selected")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Failed to hash password during user creation")
		session.Set("message", "Failed to hash password")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedpassword),
		Role:     role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		slog.Warn("Failed to create new user - possible duplicate email", "email", email, "error", err)
		session.Set("message", "Email already exists or invalid data")
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	slog.Info("New user created successfully by admin", "user_id", user.ID, "email", email, "role", role)

	session.Set("message", "User created successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
}
