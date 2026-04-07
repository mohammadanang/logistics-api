package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mohammadanang/logistics-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type packageUsecase struct {
	dbRepo    domain.PackageRepository
	cacheRepo domain.CacheRepository
}

func NewPackageUsecase(db domain.PackageRepository, cache domain.CacheRepository) domain.PackageUsecase {
	return &packageUsecase{
		dbRepo:    db,
		cacheRepo: cache,
	}
}

func (u *packageUsecase) TrackPackage(ctx context.Context, trackingNo string) (*domain.Package, error) {
	cacheKey := "track:" + trackingNo

	// 1. Coba ambil dari Cache (Redis)
	cachedData, err := u.cacheRepo.Get(ctx, cacheKey)
	if err == nil {
		var pkg domain.Package
		if err := json.Unmarshal([]byte(cachedData), &pkg); err == nil {
			log.Println("CACHE HIT: Data diambil dari Redis")
			return &pkg, nil
		}
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err) // Log error non-fatal, tetap lanjut ke DB
	}

	// 2. Jika Cache Miss, ambil dari Database (PostgreSQL)
	log.Println("CACHE MISS: Mengambil data dari PostgreSQL")
	pkg, err := u.dbRepo.FindByTrackingNo(ctx, trackingNo)
	if err != nil {
		return nil, err // Return error "not found"
	}

	// 3. Simpan hasil ke Cache dengan Time-To-Live (TTL) 5 Menit
	pkgBytes, _ := json.Marshal(pkg)
	_ = u.cacheRepo.Set(ctx, cacheKey, pkgBytes, 5*time.Minute)

	return pkg, nil
}
