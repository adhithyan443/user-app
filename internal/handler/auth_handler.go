package handler

import (
	"net/http"
	"user-app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(
	authUsecase *usecase.AuthUsecase,
) *AuthHandler {

	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) HandleSignup(ctx *gin.Context) {

	req := usecase.SignupRequest{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}

	err := h.authUsecase.Signup(req)

	if err != nil {

		ctx.HTML(
			http.StatusBadRequest,
			"signup.html",
			gin.H{
				"error": err.Error(),
				"name":  req.Name,
				"email": req.Email,
			},
		)

		return
	}

	session := sessions.Default(ctx)

	session.Set(
		"message",
		"Profile created successfully",
	)

	session.Save()

	ctx.Redirect(
		http.StatusSeeOther,
		"/login",
	)
}
