package postgres

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/google/uuid"
)

func (p *Postgres) GetPlant(ctx context.Context, id uuid.UUID) (entity.Plant, error) {
	ctx, span := tracer.Start(ctx, "postgres GetPlant")
	defer span.End()

	return entity.Plant{ID: id}, nil
}
