package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/utils"
)

type Userinfo struct {
	ID   uint64
	Role uint64
}

// AuthMiddleware : only set userinfo if token exist and valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if claims, ok := utils.ParseToken(token); ok {
			c.Set("userinfo", Userinfo{claims.ID, claims.Role})
		}
	}
}

// LoginRequired : set userinfo and abort the request if not authorized
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("userinfo"); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("userinfo")
		if user.(Userinfo).Role == 1 {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("userinfo")
		if user.(Userinfo).Role != 2 {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
