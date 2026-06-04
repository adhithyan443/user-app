package handler

import (
	"log/slog"
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

func (h *AuthHandler) HandleForgotPassword(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	email := ctx.PostForm("email")

	if email != "" {

		user, err := h.authUsecase.FindUserByEmail(
			email,
		)

		if err != nil {

			session.Set(
				"message",
				"User does not exist",
			)

			session.Save()

			ctx.Redirect(
				http.StatusSeeOther,
				"/forgotpassword",
			)

			return
		}

		session.Set(
			"reset_id",
			user.ID,
		)

		session.Save()

		ctx.HTML(
			http.StatusOK,
			"forgotpassword.html",
			gin.H{
				"account": false,
			},
		)

		return
	}

	id := session.Get("reset_id")

	userID, ok := id.(uint)

	if !ok {

		ctx.Redirect(
			http.StatusSeeOther,
			"/forgotpassword",
		)

		return
	}

	req := usecase.ResetPasswordRequest{
		UserID:          userID,
		NewPassword:     ctx.PostForm("newpassword"),
		ConfirmPassword: ctx.PostForm("confirmpassword"),
	}

	err := h.authUsecase.ResetPassword(req)

	if err != nil {

		session.Set(
			"message",
			err.Error(),
		)

		session.Save()

		ctx.Redirect(
			http.StatusSeeOther,
			"/forgotpassword",
		)

		return
	}

	session.Delete("reset_id")

	session.Set(
		"message",
		"Password updated successfully",
	)

	session.Save()

	ctx.Redirect(
		http.StatusSeeOther,
		"/login",
	)
}

func (h *AuthHandler) Logout(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	session.Clear()

	err := session.Save()

	if err != nil {

		ctx.String(
			http.StatusInternalServerError,
			"Failed to logout",
		)

		return
	}

	ctx.Redirect(
		http.StatusSeeOther,
		"/login",
	)
}
func (h *AuthHandler) JWTProfile(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get("user_id")
	email, _ := ctx.Get("email")
	role, _ := ctx.Get("role")

	slog.Info( 
		"JWT protected route accessed",
		"user_id", userID,
		"email", email,
		"role", role,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "JWT authentication successful",
			"user_id": userID,
			"email":   email,
			"role":    role,
		},
	)
}
