package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) ShowSignupPage(
	ctx *gin.Context,
) {
	ctx.HTML(
		http.StatusOK,
		"signup.html",
		nil,
	)
}

func (h *PageHandler) ShowLoginPage(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	msg := session.Get("message")

	session.Delete("message")
	session.Save()

	if session.Get("user_id") != nil {

		if session.Get("role") == "admin" {

			ctx.Redirect(
				http.StatusSeeOther,
				"/admin",
			)

		} else {

			ctx.Redirect(
				http.StatusSeeOther,
				"/home",
			)

		}

		return
	}

	ctx.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"message": msg,
		},
	)
}

func (h *PageHandler) ShowChangePasswordPage(
	ctx *gin.Context,
) {

	session := sessions.Default(ctx)

	msg := session.Get("message")
	errmsg := session.Get("error")

	session.Delete("message")
	session.Delete("error")

	session.Save()

	ctx.HTML(
		http.StatusOK,
		"user_changepassword.html",
		gin.H{
			"message": msg,
			"error":   errmsg,
		},
	)
}