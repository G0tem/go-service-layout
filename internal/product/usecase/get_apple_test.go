package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/internal/product/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetApple_Success - тест успешного получения apple из redis
func TestGetApple_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	input := dto.GetAppleInput{ID: uuid.New()}

	expectedApple := entity.Apple{
		ID:     input.ID,
		Name:   "Test Apple",
		Status: entity.StatusNew,
	}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockRedis.On("GetApple", ctx, input.ID).Return(expectedApple, nil)

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	// Act
	output, err := uc.GetApple(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedApple.ID, output.ID)
	assert.Equal(t, expectedApple.Name, output.Name)
	assert.Equal(t, expectedApple.Status, output.Status)
	mockRedis.AssertExpectations(t)
}

// TestGetApple_RedisNotFound - тест когда apple не найден в redis
func TestGetApple_RedisNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.GetAppleInput{ID: uuid.New()}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockRedis.On("GetApple", ctx, input.ID).Return(entity.Apple{}, entity.ErrNotFound)

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.GetApple(ctx, input)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, entity.ErrNotFound))
	assert.Contains(t, err.Error(), "u.redis.GetApple")
	assert.Equal(t, dto.GetAppleOutput{}, output)
	mockRedis.AssertExpectations(t)
}

// TestGetApple_RedisConnectionError - тест ошибки соединения с redis
func TestGetApple_RedisConnectionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.GetAppleInput{ID: uuid.New()}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	mockRedis.On("GetApple", ctx, input.ID).Return(entity.Apple{}, errors.New("redis connection failed"))

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.GetApple(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "u.redis.GetApple")
	assert.Contains(t, err.Error(), "redis connection failed")
	assert.Equal(t, dto.GetAppleOutput{}, output)
	mockRedis.AssertExpectations(t)
}

// TestGetApple_InvalidUUID - тест с невалидным UUID (если бы была валидация)
func TestGetApple_InvalidUUID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := dto.GetAppleInput{ID: uuid.Nil}

	mockPG := new(MockPostgres)
	mockKafka := new(MockKafka)
	mockRedis := new(MockRedis)
	mockTM := new(MockTokenManager)

	// Redis вернет NotFound для nil UUID
	mockRedis.On("GetApple", ctx, uuid.Nil).Return(entity.Apple{}, entity.ErrNotFound)

	uc := usecase.New(mockPG, mockKafka, mockRedis, mockTM)

	output, err := uc.GetApple(ctx, input)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, entity.ErrNotFound))
	assert.Equal(t, dto.GetAppleOutput{}, output)
}
