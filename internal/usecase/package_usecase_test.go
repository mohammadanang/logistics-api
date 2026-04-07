package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadanang/logistics-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository adalah tiruan dari PackageRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByTrackingNo(ctx context.Context, trackingNo string) (*domain.Package, error) {
	args := m.Called(ctx, trackingNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Package), args.Error(1)
}

// MockCache adalah tiruan dari CacheRepository
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return m.Called(ctx, key, value, exp).Error(0)
}

// TEST CASE: Berhasil tracking saat cache kosong (Cache Miss)
func TestTrackPackage_Success_CacheMiss(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache)
	usecase := NewPackageUsecase(mockRepo, mockCache)

	trackingNo := "RESI-123"
	expectedPkg := &domain.Package{TrackingNo: trackingNo, Status: "MANIFESTED"}

	// Ekspektasi: Cache return kosong (Error redis.Nil disimulasikan sebagai empty string/err)
	mockCache.On("Get", mock.Anything, "track:"+trackingNo).Return("", errors.New("redis: nil"))

	// Ekspektasi: Database return data paket
	mockRepo.On("FindByTrackingNo", mock.Anything, trackingNo).Return(expectedPkg, nil)

	// Ekspektasi: Hasil dari DB disimpan ke Cache (Set)
	mockCache.On("Set", mock.Anything, "track:"+trackingNo, mock.Anything, mock.Anything).Return(nil)

	// Eksekusi
	result, err := usecase.TrackPackage(context.Background(), trackingNo)

	// Validasi
	assert.NoError(t, err)
	assert.Equal(t, expectedPkg.Status, result.Status)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
