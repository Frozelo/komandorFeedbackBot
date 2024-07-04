package repository

import (
	"context"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
	"github.com/jackc/pgx/v4"
)

type SurveyRepository struct {
	db *pgx.Conn
}

func NewSurveyRepository(db *pgx.Conn) *SurveyRepository {
	return &SurveyRepository{db: db}
}

func (r *SurveyRepository) CreateSurvey(survey entity.Survey) (*entity.Survey, error) {
	query := `INSERT INTO surveys (user_tg_id, question, answer) VALUES ($1, $2, $3) RETURNING id, user_tg_id, question, answer`
	row := r.db.QueryRow(context.Background(), query, survey.UserId, survey.Question, survey.Answer)

	var newSurvey entity.Survey
	err := row.Scan(&newSurvey.Id, &newSurvey.UserId, &newSurvey.Question, &newSurvey.Answer)
	if err != nil {
		return nil, err
	}

	return &newSurvey, nil
}

func (r *SurveyRepository) GetSurveyResults(userId int) ([]entity.Survey, error) {
	query := `SELECT id, user_tg_id, question, answer FROM surveys WHERE user_tg_id = $1`
	rows, err := r.db.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []entity.Survey
	for rows.Next() {
		var survey entity.Survey
		if err := rows.Scan(&survey.Id, &survey.UserId, &survey.Question, &survey.Answer); err != nil {
			return nil, err
		}
		surveys = append(surveys, survey)
	}

	return surveys, nil
}
