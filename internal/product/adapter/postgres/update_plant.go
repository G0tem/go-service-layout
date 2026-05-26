package postgres

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/doug-martin/goqu/v9"
)

func (p *Postgres) UpdatePlant(ctx context.Context, pl entity.Plant) (err error) {
	ctx, span := tracer.Start(ctx, "postgres UpdatePlant")
	defer span.End()

	dataset := goqu.Update("plant").
		Set(goqu.Record{
			"name":   string(pl.Name),
			"status": pl.Status,
		}).
		Where(goqu.Ex{"id": pl.ID})

	sql, args, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("dataset.ToSQL: %w", err)
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("p.pool.Exec: %w", err)
	}

	return nil
}
