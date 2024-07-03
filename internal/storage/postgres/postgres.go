package postgres_storage

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type PostgresStorage struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  int

	Pool *pgxpool.Pool
}
