package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammadanang/logistics-api/internal/usecase"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

// Injeksi PaymentUsecase ke dalam Handler
func NewPaymentHandler(r *gin.RouterGroup, pu *usecase.PaymentUsecase) {
	handler := &PaymentHandler{paymentUsecase: pu}

	// Kita letakkan di dalam rute yang terproteksi (harus login)
	// Atau jika pelanggan yang membuat, rute ini bisa diletakkan di rute publik
	r.POST("/payments/invoice", handler.CreateInvoice)
}

func (h *PaymentHandler) CreateInvoice(c *gin.Context) {
	// 1. Definisikan struktur payload yang diharapkan dari Postman/Frontend
	var req struct {
		TrackingNo    string  `json:"tracking_no" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		CustomerEmail string  `json:"customer_email" binding:"required,email"`
	}

	// 2. Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Format data tidak valid: " + err.Error(),
		})
		return
	}

	// 3. Panggil Usecase untuk generate link Xendit
	invoiceURL, err := h.paymentUsecase.CreateShippingInvoice(
		c.Request.Context(),
		req.TrackingNo,
		req.Amount,
		req.CustomerEmail,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat tagihan pembayaran",
			"details": err.Error(),
		})
		return
	}

	// 4. Kembalikan URL Invoice ke client
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Invoice berhasil dibuat",
		"invoice_url": invoiceURL,
	})
}
