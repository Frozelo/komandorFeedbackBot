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
	query := `INSERT INTO surveys (user_id) VALUES ($1) RETURNING id, user_id, avg_score`
	row := r.db.QueryRow(context.Background(), query, survey.UserId)

	var newSurvey entity.Survey
	err := row.Scan(&newSurvey.Id, &newSurvey.UserId, &newSurvey.AvgScore)
	if err != nil {
		return nil, err
	}

	return &newSurvey, nil
}

func (r *SurveyRepository) UpdateQuestionAnswer(questionId int, answer int) error {
	query := `UPDATE answers SET answer = $1 WHERE id = $2`
	_, err := r.db.Exec(context.Background(), query, answer, questionId)
	return err
}

func (r *SurveyRepository) GetSurveyQuestion(questionId int) (*entity.Question, error) {
	query := `SELECT id, text, category_id FROM questions WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, questionId)

	var question entity.Question
	err := row.Scan(&question.Id, &question.Text, &question.CategoryId)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *SurveyRepository) GetSurvey(surveyId int) (*entity.Survey, error) {
	query := `SELECT id, user_id, avg_score FROM surveys WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, surveyId)

	var survey entity.Survey
	err := row.Scan(&survey.Id, &survey.UserId, &survey.AvgScore)
	if err != nil {
		return nil, err
	}

	return &survey, nil
}

func (r *SurveyRepository) CalculateAverageScore(surveyId int) (float64, error) {
	query := `SELECT AVG(answer) FROM answers WHERE survey_id = $1`
	row := r.db.QueryRow(context.Background(), query, surveyId)

	var avgScore float64
	err := row.Scan(&avgScore)
	if err != nil {
		return 0, err
	}

	return avgScore, nil
}

func (r *SurveyRepository) SaveAvgScore(surveyId int, avgScore float64) error {
	query := `UPDATE surveys SET avg_score = $1 WHERE id = $2`
	_, err := r.db.Exec(context.Background(), query, avgScore, surveyId)
	return err
}

func (r *SurveyRepository) GetCategories() ([]entity.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *SurveyRepository) GetQuestionsByCategory(categoryID int) ([]entity.Question, error) {
	query := `SELECT id, text, category_id FROM questions WHERE category_id = $1`
	rows, err := r.db.Query(context.Background(), query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		err := rows.Scan(&question.Id, &question.Text, &question.CategoryId)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (r *SurveyRepository) SaveAnswer(answer entity.Answer) error {
	query := `INSERT INTO answers (survey_id, question_id, answer) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(context.Background(), query, answer.SurveyID, answer.QuestionID, answer.Answer)
	return err
}
