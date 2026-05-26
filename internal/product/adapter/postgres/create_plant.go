package postgres

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/doug-martin/goqu/v9"
)

func (p *Postgres) CreatePlant(ctx context.Context, pl entity.Plant) (err error) {
	ctx, span := tracer.Start(ctx, "postgres CreatePlant")
	defer span.End()

	dataset := goqu.Insert("plant").
		Rows(goqu.Record{
			"id":     pl.ID,
			"name":   pl.Name,
			"status": pl.Status,
		})

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("dataset.ToSQL: %w", err)
	}

	_, err = p.pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}
