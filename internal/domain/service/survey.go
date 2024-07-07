package service

import (
	"sync"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
)

type SurveyRepository interface {
	CreateSurvey(survey entity.Survey) (*entity.Survey, error)
	UpdateQuestionAnswer(questionId int, answer int) error
	GetSurveyQuestion(questionId int) (*entity.Question, error)
	GetSurvey(surveyId int) (*entity.Survey, error)
	CalculateAverageScore(surveyId int) (float64, error)
	SaveAvgScore(surveyId int, avgScore float64) error
	GetCategories() ([]entity.Category, error)
	GetQuestionsByCategory(categoryID int) ([]entity.Question, error)
	SaveAnswer(answer entity.Answer) error
}

type SurveyService struct {
	repo SurveyRepository
	mu   sync.Mutex
}

func NewSurveyService(repo SurveyRepository) *SurveyService {
	return &SurveyService{repo: repo}
}

func (s *SurveyService) CreateSurvey(survey entity.Survey) (*entity.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.CreateSurvey(survey)
}

func (s *SurveyService) UpdateQuestionAnswer(questionId int, answer int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.UpdateQuestionAnswer(questionId, answer)
}

func (s *SurveyService) GetSurveyQuestion(questionId int) (*entity.Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.GetSurveyQuestion(questionId)
}

func (s *SurveyService) GetSurvey(surveyId int) (*entity.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.GetSurvey(surveyId)
}

func (s *SurveyService) CalculateAverageScore(surveyId int) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.CalculateAverageScore(surveyId)
}

func (s *SurveyService) SaveAvgScore(surveyId int, avgScore float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.SaveAvgScore(surveyId, avgScore)
}

func (s *SurveyService) GetCategories() ([]entity.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.GetCategories()
}

func (s *SurveyService) GetQuestionsByCategory(categoryID int) ([]entity.Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.GetQuestionsByCategory(categoryID)
}

func (s *SurveyService) SaveAnswer(answer entity.Answer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.SaveAnswer(answer)
}
