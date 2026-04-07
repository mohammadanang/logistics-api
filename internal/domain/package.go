package domain

import (
	"context"
	"time"
)

// Package mewakili entitas paket dalam sistem logistik
type Package struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TrackingNo   string    `json:"tracking_no" gorm:"uniqueIndex;not null"`
	SenderName   string    `json:"sender_name"`
	ReceiverName string    `json:"receiver_name"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Status       string    `json:"status"` // e.g., "MANIFESTED", "ON_TRANSIT", "DELIVERED"
	UpdatedAt    time.Time `json:"updated_at"`
}

// PackageRepository adalah kontrak/interface untuk akses database
type PackageRepository interface {
	FindByTrackingNo(ctx context.Context, trackingNo string) (*Package, error)
	// Kita akan tambahkan BatchUpdate nanti
}

// CacheRepository adalah kontrak/interface untuk akses Redis
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

// PackageUsecase adalah kontrak/interface untuk business logic
type PackageUsecase interface {
	TrackPackage(ctx context.Context, trackingNo string) (*Package, error)
}
