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