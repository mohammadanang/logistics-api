package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohammadanang/logistics-api/pkg/paseto"
)

func AuthMiddleware(tokenMaker *paseto.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is not provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Ambil data dari payload dan simpan di context Gin agar bisa dipakai oleh handler
		userID, _ := payload.GetString("user_id")
		role, _ := payload.GetString("role")

		c.Set("x-user-id", userID)
		c.Set("x-user-role", role)

		c.Next() // Lanjut ke handler berikutnya
	}
}
