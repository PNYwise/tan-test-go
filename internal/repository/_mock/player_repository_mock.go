package mock_test

import (
	"tan-test-go/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) CreateBatch(players *[]domain.Player) error {
	args := m.Called(players)
	return args.Error(0)
}

func (m *MockPlayerRepository) GetPlayers() (*[]domain.Player, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Player), args.Error(1)
}
