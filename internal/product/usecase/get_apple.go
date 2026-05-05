package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
)

func (u *UseCase) GetApple(ctx context.Context, input dto.GetAppleInput) (dto.GetAppleOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase GetApple")
	defer span.End()

	var output dto.GetAppleOutput

	apple, err := u.redis.GetApple(ctx, input.ID)
	if err != nil {
		return output, fmt.Errorf("u.redis.GetApple: %w", err)
	}

	return dto.GetAppleOutput{
		ID:     apple.ID,
		Name:   apple.Name,
		Status: apple.Status,
	}, err
}
