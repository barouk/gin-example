package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AccessControlMiddleware(allowedRoles []int) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := getUserRole(c)
		if !isRoleAllowed(role.(int), allowedRoles) || err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

func getUserRole(c *gin.Context) (interface{}, error) {
	role, exists := c.Get("GroupID")
	if !exists {
		return "", errors.New("role dont exists ")
	}
	return role.(interface{}), nil
}

func isRoleAllowed(role int, allowedRoles []int) bool {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
