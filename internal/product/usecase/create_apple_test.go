package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/internal/product/usecase"
	jwt "github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostgres - моки для postgres адаптера
type MockPostgres struct {
	mock.Mock
}

func (m *MockPostgres) CreateApple(ctx context.Context, a entity.Apple) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockPostgres) GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Apple), args.Error(1)
}

func (m *MockPostgres) CreatePineApple(ctx context.Context, p entity.PineApple) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

// MockKafka - моки для kafka адаптера
type MockKafka struct {
	mock.Mock
}

func (m *MockKafka) CreateEvent(ctx context.Context, e entity.CreateEvent) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

// MockRedis - моки для redis адаптера
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Apple), args.Error(1)
}

func (m *MockRedis) PutApple(ctx context.Context, a entity.Apple) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

// MockTokenManager - моки для JWT менеджера
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateToken(ctx context.Context, claims jwt.AuthClaims) (string, error) {
	args := m.Called(ctx, claims)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) ValidateToken(ctx context.Context, token string) (*jwt.AuthClaims, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*jwt.AuthClaims), args.Error(1)
}

// TestCreateApple_Success - тест успешного создания apple
func TestCreateApple_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	input := dto.CreateAppleInput{Name: "Test Apple"}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	// Ожидаем вызовы в правильном порядке
	mockPG.On("CreateApple", mock.Anything, mock.MatchedBy(func(a entity.Apple) bool {
		return a.Name == input.Name && a.Status == entity.StatusNew
	})).Return(nil)

	mockRedis.On("PutApple", mock.Anything, mock.MatchedBy(func(a entity.Apple) bool {
		return a.Name == input.Name && a.Status == entity.StatusNew
	})).Return(nil)

	mockKafka.On("CreateEvent", mock.Anything, mock.MatchedBy(func(e entity.CreateEvent) bool {
		return e.Name == input.Name
	})).Return(nil)

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	// Act
	output, err := uc.CreateApple(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, output.ID)
	mockPG.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

// TestCreateApple_PostgresError - тест ошибки при записи в postgres
func TestCreateApple_PostgresError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.CreateAppleInput{Name: "Test Apple"}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockPG.On("CreateApple", mock.Anything, mock.Anything).Return(errors.New("db connection failed"))

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.CreateApple(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "u.postgres.CreateApple")
	assert.Equal(t, dto.CreateAppleOutput{}, output)
	mockPG.AssertExpectations(t)
	mockRedis.AssertNotCalled(t, "PutApple", mock.Anything, mock.Anything)
	mockKafka.AssertNotCalled(t, "CreateEvent", mock.Anything, mock.Anything)
}

// TestCreateApple_RedisError - тест ошибки при записи в redis
func TestCreateApple_RedisError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.CreateAppleInput{Name: "Test Apple"}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockPG.On("CreateApple", mock.Anything, mock.MatchedBy(func(a entity.Apple) bool {
		return a.Name == input.Name
	})).Return(nil)

	mockRedis.On("PutApple", mock.Anything, mock.Anything).Return(errors.New("redis connection failed"))

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.CreateApple(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "u.redis.PutApple")
	assert.Equal(t, dto.CreateAppleOutput{}, output)
	mockPG.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockKafka.AssertNotCalled(t, "CreateEvent", mock.Anything, mock.Anything)
}

// TestCreateApple_KafkaError - тест ошибки при отправке в kafka
func TestCreateApple_KafkaError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.CreateAppleInput{Name: "Test Apple"}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockPG.On("CreateApple", mock.Anything, mock.MatchedBy(func(a entity.Apple) bool {
		return a.Name == input.Name
	})).Return(nil)

	mockRedis.On("PutApple", mock.Anything, mock.MatchedBy(func(a entity.Apple) bool {
		return a.Name == input.Name
	})).Return(nil)

	mockKafka.On("CreateEvent", mock.Anything, mock.Anything).Return(errors.New("kafka unavailable"))

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.CreateApple(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "u.kafka.CreateEvent")
	assert.Equal(t, dto.CreateAppleOutput{}, output)
	mockPG.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}
