package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: p,
	}
}
