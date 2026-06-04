package handler

import (
	"net/http"
	"user-app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(
	userUsecase *usecase.UserUsecase,
) *UserHandler {

	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) ShowProfilePage(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	id := session.Get("user_id")

	msg := session.Get("message")

	session.Delete("message")
	session.Save()

	userID, ok := id.(uint)

	if !ok {
		ctx.String(
			http.StatusUnauthorized,
			"invalid session",
		)
		return
	}

	user, err := h.userUsecase.GetProfile(userID)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"user not found",
		)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"user_profile.html",
		gin.H{
			"user":    user,
			"message": msg,
		},
	)
}

func (h *UserHandler) UpdateUserProfile(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	id := session.Get("user_id")

	userID, ok := id.(uint)

	if !ok {
		ctx.Redirect(
			http.StatusSeeOther,
			"/login",
		)
		return
	}

	req := usecase.UpdateProfileRequest{
		ID:    userID,
		Name:  ctx.PostForm("name"),
		Email: ctx.PostForm("email"),
	}

	err := h.userUsecase.UpdateProfile(req)

	if err != nil {

		session.Set(
			"message",
			err.Error(),
		)

		session.Save()

		ctx.Redirect(
			http.StatusSeeOther,
			"/profile",
		)

		return
	}

	session.Set(
		"message",
		"Profile updated successfully",
	)

	session.Set(
		"name",
		req.Name,
	)

	session.Set(
		"email",
		req.Email,
	)

	session.Save()

	ctx.Redirect(
		http.StatusSeeOther,
		"/profile",
	)
}
