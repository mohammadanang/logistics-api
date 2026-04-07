package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadanang/logistics-api/pkg/paseto"
)

type AuthHandler struct {
	tokenMaker *paseto.TokenMaker
}

func NewAuthHandler(r *gin.RouterGroup, tokenMaker *paseto.TokenMaker) {
	handler := &AuthHandler{tokenMaker: tokenMaker}
	r.POST("/login", handler.Login)
}

func (h *AuthHandler) Login(c *gin.Context) {
	// SIMULASI: Di production, Anda harus query DB untuk cek username & password dengan bcrypt
	// Untuk portofolio, kita asumsikan login selalu sukses

	token, err := h.tokenMaker.CreateToken("courier-77", "COURIER", 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}
