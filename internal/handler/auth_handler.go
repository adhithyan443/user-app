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

func (h *AuthHandler) HandleLogin(ctx *gin.Context) {

	req := usecase.LoginRequest{
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}

	res, err := h.authUsecase.Login(req)

	if err != nil {

		ctx.HTML(
			http.StatusUnauthorized,
			"login.html",
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	session := sessions.Default(ctx)

	session.Set("user_id", res.User.ID)
	session.Set("email", res.User.Email)
	session.Set("name", res.User.Name)
	session.Set("role", res.User.Role)
	session.Set("token", res.Token)

	if err := session.Save(); err != nil {
		ctx.HTML(
			http.StatusInternalServerError,
			"login.html",
			gin.H{
				"error": "failed to save session",
			},
		)
		return
	}

	if res.User.Role == "admin" {
		ctx.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/home")
}