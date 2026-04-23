package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func AdminRequired() gin.HandlerFunc{
	return func(ctx *gin.Context) {

		session:=sessions.Default(ctx)

		role:=session.Get("role")

		if role != "admin"{
			// ctx.HTML(http.StatusForbidden, "error.html",gin.H{
			// 	"error":"Access denied. Admin only.",
			// })
			ctx.Redirect(http.StatusFound, "/home")
			ctx.Abort()
			return 
		}

		ctx.Next()
	}
}