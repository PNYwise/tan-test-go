package mock_test

import (
	"tan-test-go/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockGeolocationRepository struct {
	mock.Mock
}

func (m *MockGeolocationRepository) CreateBatch(geolocations *[]domain.Geolocation) error {
	args := m.Called(geolocations)
	return args.Error(0)
}

func (m *MockGeolocationRepository) GetGeolocations() (*[]domain.Geolocation, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Geolocation), args.Error(1)
}
