package repository

import (
	"context"
	"errors"

	"github.com/mohammadanang/logistics-api/internal/domain"
	"gorm.io/gorm"
)

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.PackageRepository {
	// AutoMigrate hanya untuk mempermudah setup awal portofolio
	_ = db.AutoMigrate(&domain.Package{})
	return &postgresRepo{db: db}
}

func (p *postgresRepo) FindByTrackingNo(ctx context.Context, trackingNo string) (*domain.Package, error) {
	var pkg domain.Package
	// Menggunakan WithContext sangat penting untuk mencegah query menggantung (timeout)
	err := p.db.WithContext(ctx).Where("tracking_no = ?", trackingNo).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("package not found")
		}
		return nil, err
	}
	return &pkg, nil
}
