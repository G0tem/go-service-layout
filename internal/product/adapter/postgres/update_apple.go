package postgres

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/doug-martin/goqu/v9"
)

func (p *Postgres) UpdateApple(ctx context.Context, a entity.Apple) (err error) {
	ctx, span := tracer.Start(ctx, "postgres UpdateApple")
	defer span.End()

	dataset := goqu.Update("apple").
		Set(goqu.Record{
			"name":   a.Name,
			"status": a.Status,
		}).
		Where(goqu.Ex{"id": a.ID})

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
