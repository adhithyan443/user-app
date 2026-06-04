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
