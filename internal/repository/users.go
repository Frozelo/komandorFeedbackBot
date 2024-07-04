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

func (r *UserRepository) GetUserByTgId(tgId int) (*entity.User, error) {
	query := `SELECT id, tg_id, username FROM users WHERE tg_id = $1`
	row := r.db.QueryRow(context.Background(), query, tgId)

	var user entity.User
	err := row.Scan(&user.Id, &user.TgId, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil

}

func (r *UserRepository) CreateTgUser(user entity.User) (*entity.User, error) {
	query := `INSERT INTO users (tg_id, username) VALUES ($1, $2) RETURNING id, tg_id, username`
	row := r.db.QueryRow(context.Background(), query, user.TgId, user.Username)

	var newUser entity.User
	err := row.Scan(&newUser.Id, &newUser.TgId, &newUser.Username)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}
