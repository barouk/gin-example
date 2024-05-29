package middleware

import (
	"gin-example/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AthorizationJWT(jwtService jwt.Jwt) gin.HandlerFunc {
	return func(context *gin.Context) {
		authToken := context.GetHeader("Authorization")
		if authToken == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token", "status": false})
			context.Abort()
			return
		}
		user, err := jwtService.ValidateToken(authToken)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "status": false})
			context.Abort()
			return
		}
		context.Set("username", user.Username)
		context.Set("id", user.Id)
		return
	}
}
