package postgres_storage

import (
	"context"
	"fmt"

	"github.com/Frozelo/komandorFeedbackBot/internal/bot"
	"github.com/jackc/pgx/v4"

	_ "github.com/lib/pq"
)

func NewDb(cfg *bot.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DatabaseName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
