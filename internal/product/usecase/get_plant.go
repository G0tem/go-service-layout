package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
)

func (u *UseCase) GetPlant(ctx context.Context, input dto.GetPlantInput) (dto.GetPlantOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase GetPlant")
	defer span.End()

	var output dto.GetPlantOutput

	plant, err := u.redis.GetPlant(ctx, input.ID)
	if err != nil {
		return output, fmt.Errorf("u.redis.GetPlant: %w", err)
	}

	return dto.GetPlantOutput{
		ID:     plant.ID,
		Name:   string(plant.Name),
		Status: int(plant.Status),
	}, err
}
