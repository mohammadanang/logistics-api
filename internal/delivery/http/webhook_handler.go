package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	callbackToken string
}

func NewWebhookHandler(r *gin.RouterGroup, callbackToken string) {
	handler := &WebhookHandler{callbackToken: callbackToken}
	r.POST("/webhooks/xendit/invoice", handler.HandleInvoiceCallback)
}

func (h *WebhookHandler) HandleInvoiceCallback(c *gin.Context) {
	// 1. VERIFIKASI KEAMANAN (Wajib)
	reqToken := c.GetHeader("x-callback-token")
	if reqToken != h.callbackToken {
		log.Println("SECURITY ALERT: Webhook token tidak valid!")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized webhook"})
		return
	}

	// 2. Parse Payload Invoice Xendit
	var payload struct {
		ID             string  `json:"id"`
		ExternalID     string  `json:"external_id"` // Ini adalah Nomor Resi kita
		Status         string  `json:"status"`      // "PAID", "EXPIRED", dll
		PaidAmount     float64 `json:"paid_amount"`
		PaymentMethod  string  `json:"payment_method"`
		PaymentChannel string  `json:"payment_channel"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Gagal parsing JSON webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak sesuai"})
		return
	}

	// 3. Eksekusi Logika Bisnis (Update Database)
	if payload.Status == "PAID" || payload.Status == "SETTLED" {
		log.Printf("SUKSES: Pembayaran diterima untuk Resi %s sebesar %.2f via %s",
			payload.ExternalID, payload.PaidAmount, payload.PaymentChannel)

		// TODO: Panggil usecase untuk mengubah status resi di database (misal dari "MENUNGGU_PEMBAYARAN" menjadi "DIPROSES")
	} else if payload.Status == "EXPIRED" {
		log.Printf("EXPIRED: Tagihan untuk Resi %s telah kedaluwarsa", payload.ExternalID)
		// TODO: Panggil usecase untuk membatalkan pengiriman
	}

	// 4. Balas Xendit dengan 200 OK agar mereka tahu kita sudah menerimanya
	c.JSON(http.StatusOK, gin.H{"message": "Webhook diterima dengan baik"})
}
