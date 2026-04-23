package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		session:=sessions.Default(ctx)

		userID:=session.Get("user_id")

		if userID == nil{
			ctx.Redirect(http.StatusFound,"/login")
			ctx.Abort()
			return 
		}

		ctx.Next()
	}
}