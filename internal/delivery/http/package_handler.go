package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammadanang/logistics-api/internal/domain"
)

type PackageHandler struct {
	usecase domain.PackageUsecase
}

func NewPackageHandler(r *gin.Engine, us domain.PackageUsecase) {
	handler := &PackageHandler{usecase: us}

	// Grouping versi API
	v1 := r.Group("/api/v1")
	{
		v1.GET("/track/:resi", handler.TrackPackage)
	}
}

func (h *PackageHandler) TrackPackage(c *gin.Context) {
	resi := c.Param("resi")

	// Panggil usecase dengan context dari Gin
	pkg, err := h.usecase.TrackPackage(c.Request.Context(), resi)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Package found",
		"data":    pkg,
	})
}
