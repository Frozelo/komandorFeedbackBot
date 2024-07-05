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
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `INSERT INTO surveys (user_id) VALUES ($1) RETURNING id`
	row := tx.QueryRow(context.Background(), query, survey.UserId)

	var surveyId int
	if err := row.Scan(&surveyId); err != nil {
		return nil, err
	}
	survey.Id = surveyId

	for i, question := range survey.Questions {
		query := `INSERT INTO questions (survey_id, text) VALUES ($1, $2) RETURNING id`
		row := tx.QueryRow(context.Background(), query, surveyId, question.Text)
		var questionId int
		if err := row.Scan(&questionId); err != nil {
			return nil, err
		}
		survey.Questions[i].Id = questionId
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return &survey, nil
}

func (r *SurveyRepository) UpdateQuestionAnswer(questionId int, answer int) error {
	query := `UPDATE questions SET answer = $1 WHERE id = $2`
	_, err := r.db.Exec(context.Background(), query, answer, questionId)
	return err
}

func (r *SurveyRepository) GetSurveyQuestion(questionId int) (*entity.Question, error) {
	query := `SELECT id, survey_id, text, answer FROM questions WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, questionId)

	var question entity.Question
	err := row.Scan(&question.Id, &question.SurveyId, &question.Text, &question.Answer)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *SurveyRepository) GetSurvey(surveyId int) (*entity.Survey, error) {
	query := `SELECT id, user_id FROM surveys WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, surveyId)

	var survey entity.Survey
	err := row.Scan(&survey.Id, &survey.UserId)
	if err != nil {
		return nil, err
	}

	query = `SELECT id, survey_id, text, answer FROM questions WHERE survey_id = $1`
	rows, err := r.db.Query(context.Background(), query, surveyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		if err := rows.Scan(&question.Id, &question.SurveyId, &question.Text, &question.Answer); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	survey.Questions = questions

	return &survey, nil
}

func (r *SurveyRepository) CalculateAverageScore(surveyId int) (float64, error) {
	query := `SELECT AVG(answer) FROM questions WHERE survey_id = $1`
	row := r.db.QueryRow(context.Background(), query, surveyId)

	var averageScore float64
	err := row.Scan(&averageScore)
	if err != nil {
		return 0, err
	}

	return averageScore, nil
}

func (r *SurveyRepository) SaveAvgScore(surveyId int, avgScore float64) error {
	query := `UPDATE surveys SET avg_score  =  $1 WHERE id  =  $2`
	_, err := r.db.Exec(context.Background(), query, avgScore, surveyId)
	return err
}
