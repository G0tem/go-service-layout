package postgres

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
)

func (p *Postgres) CreatePlant(ctx context.Context, _ entity.Plant) (err error) {
	ctx, span := tracer.Start(ctx, "postgres CreatePlant")
	defer span.End()

	return nil
}
