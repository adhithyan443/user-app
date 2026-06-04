package handler

import (
	"net/http"

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
