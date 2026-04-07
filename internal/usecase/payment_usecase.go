package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	xendit "github.com/xendit/xendit-go/v4"
	invoice "github.com/xendit/xendit-go/v4/invoice"
)

type PaymentUsecase struct {
	xenditClient *xendit.APIClient
}

// Injeksi client Xendit dari main.go
func NewPaymentUsecase(apiKey string) *PaymentUsecase {
	client := xendit.NewClient(apiKey)
	return &PaymentUsecase{xenditClient: client}
}

// CreateShippingInvoice men-generate link pembayaran ongkir
func (u *PaymentUsecase) CreateShippingInvoice(ctx context.Context, trackingNo string, amount float64, customerEmail string) (string, error) {
	// Menentukan batas waktu pembayaran (misal: 24 jam)
	duration := 24 * time.Hour
	invoiceDuration := int32(duration.Seconds())
	invoiceDurationStr := fmt.Sprintf("%d", invoiceDuration)

	// Membentuk Request sesuai standar Xendit v4
	req := *invoice.NewCreateInvoiceRequest(
		trackingNo, // External ID (harus unik, kita pakai nomor resi)
		amount,     // Jumlah pembayaran ongkir
	)
	req.PayerEmail = &customerEmail
	description := "Pembayaran Ongkos Kirim untuk Resi: " + trackingNo
	req.Description = &description
	req.InvoiceDuration = &invoiceDurationStr

	// Eksekusi pemanggilan API ke Xendit
	resp, _, err := u.xenditClient.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(req).Execute()
	if err != nil {
		log.Printf("Gagal membuat invoice Xendit: %v\n", err)
		return "", err
	}

	// Mengembalikan URL Pembayaran yang akan diberikan ke frontend/pelanggan
	return resp.InvoiceUrl, nil
}
