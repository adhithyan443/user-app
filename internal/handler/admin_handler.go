package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"user-app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUsecase *usecase.AdminUsecase
}

func NewAdminHandler(
	adminUsecase *usecase.AdminUsecase,
) *AdminHandler {

	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

func (h *AdminHandler) GetAllUsers(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	msg := session.Get("message")

	session.Delete("message")
	session.Save()

	users, err := h.adminUsecase.GetAllUsers()

	if err != nil {

		ctx.String(
			http.StatusInternalServerError,
			"Error fetching users",
		)

		return
	}

	ctx.HTML(
		http.StatusOK,
		"admin_users.html",
		gin.H{
			"users":   users,
			"message": msg,
		},
	)
}

func (h *AdminHandler) DeleteUser(
	ctx *gin.Context,
) {

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {

		slog.Warn(
			"Invalid user ID format",
			"id_param", idParam,
		)

		ctx.String(
			http.StatusBadRequest,
			"Invalid user ID",
		)

		return
	}

	err = h.adminUsecase.DeleteUser(
		uint(id),
	)

	if err != nil {

		slog.Error(
			"Failed to delete user",
			"user_id", id,
			"error", err,
		)

		ctx.String(
			http.StatusInternalServerError,
			"Delete failed",
		)

		return
	}

	session := sessions.Default(ctx)

	session.Set(
		"message",
		"User deleted successfully",
	)

	session.Save()

	ctx.Redirect(
		http.StatusFound,
		"/admin/users",
	)
}

func (h *AdminHandler) EditUserPage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	msg := session.Get("message")
	session.Delete("message")
	session.Save()

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID format", "id_param", idParam)
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.adminUsecase.GetUserByID(uint(id))
	if err != nil {
		slog.Error("Failed to fetch user for edit", "user_id", id, "error", err)
		ctx.String(http.StatusNotFound, "User not found")
		return
	}

	ctx.HTML(http.StatusOK, "edit_user.html", gin.H{
		"user":    user,
		"message": msg,
	})
}

// UpdateUserPage - POST /admin/users/update/:id
func (h *AdminHandler) UpdateUserPage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID format in update", "id_param", idParam)
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	req := usecase.UpdateUserRequest{
		ID:    uint(id),
		Name:  ctx.PostForm("name"),
		Email: ctx.PostForm("email"),
		Role:  ctx.PostForm("role"),
	}

	session := sessions.Default(ctx)

	err = h.adminUsecase.UpdateUser(req)
	if err != nil {
		slog.Warn("Validation or update failed", "user_id", id, "error", err)
		session.Set("message", err.Error())
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/edit/"+idParam)
		return
	}

	slog.Info("User updated successfully via usecase", "user_id", id)
	session.Set("message", "User updated successfully")
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/users")
}

func (h *AdminHandler) NewUserPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	data := gin.H{
		"message":    session.Get("message"),
		"form_name":  session.Get("form_name"),
		"form_email": session.Get("form_email"),
		"form_role":  session.Get("form_role"),
	}

	// Clear flash data
	session.Delete("message")
	session.Delete("form_name")
	session.Delete("form_email")
	session.Delete("form_role")
	session.Save()

	ctx.HTML(http.StatusOK, "admin_add_user.html", data)
}

// AddNewUser - POST /admin/newuser
func (h *AdminHandler) AddNewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)

	req := usecase.CreateUserRequest{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
		Role:     ctx.PostForm("role"),
	}

	// Save form data for repopulation on error
	formData := map[string]interface{}{
		"form_name":  req.Name,
		"form_email": req.Email,
		"form_role":  req.Role,
	}

	err := h.adminUsecase.CreateUser(req)
	if err != nil {
		slog.Warn("Create user failed", "error", err)
		session.Set("message", err.Error())
		for k, v := range formData {
			session.Set(k, v)
		}
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/newuser")
		return
	}

	slog.Info("New user created successfully by admin", "email", req.Email)
	session.Set("message", "User created successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
}

// ShowUserPasswordPage - GET /admin/users/updatepassword/:id
func (h *AdminHandler) ShowUserPasswordPage(ctx *gin.Context) {
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

// EditUserPasswordPage - POST /admin/users/updatepassword/:id
func (h *AdminHandler) EditUserPasswordPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID for password update", "id", idParam)
		ctx.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	newPassword := ctx.PostForm("newpassword")
	confirmPassword := ctx.PostForm("confirmpassword")

	if newPassword != confirmPassword {
		session.Set("message", "Passwords do not match")
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+idParam)
		return
	}

	err = h.adminUsecase.UpdateUserPassword(uint(id), newPassword)
	if err != nil {
		slog.Warn("Password update failed", "user_id", id, "error", err)
		session.Set("message", err.Error())
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/admin/users/updatepassword/"+idParam)
		return
	}

	slog.Info("Admin updated user password successfully", "user_id", id)
	session.Set("message", "Password updated successfully")
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/admin/users")
}