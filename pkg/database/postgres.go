package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConn(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Konfigurasi Connection Pool (Nilai Jual Eksekutif)
	sqlDB.SetMaxIdleConns(10)           // Jumlah minimal koneksi standby
	sqlDB.SetMaxOpenConns(100)          // Batas maksimal koneksi simultan (mencegah DB crash)
	sqlDB.SetConnMaxLifetime(time.Hour) // Umur maksimal koneksi sebelum di-refresh

	log.Println("PostgreSQL connected successfully with Connection Pooling")
	return db
}
