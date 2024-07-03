package repository

import (
	"context"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
	"github.com/jackc/pgx/v4"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateTgUser(user entity.User) (entity.User, error) {
	query := `INSERT INTO users (tg_id, username, joined_at) VALUES ($1, $2, $3) RETURNING *`
	row := r.db.QueryRow(context.Background(), query, user.TgId, user.Username, user.JoinedAt)

	var newUser entity.User
	err := row.Scan(&newUser.Id, &newUser.TgId, &newUser.Username, &newUser.JoinedAt)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}