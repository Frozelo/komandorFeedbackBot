package service

import (
	"sync"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
)

type SurveyRepository interface {
	CreateSurvey(survey entity.Survey) (*entity.Survey, error)
	GetSurveyResults(userId int) ([]entity.Survey, error)
}

type SurveyService struct {
	mu   sync.Mutex
	repo SurveyRepository
}

func NewSurveyService(repo SurveyRepository) *SurveyService {
	return &SurveyService{repo: repo}
}

func (s *SurveyService) CreateSurvey(survey entity.Survey) (*entity.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.CreateSurvey(survey)
}

func (s *SurveyService) GetSurveyResults(userId int) ([]entity.Survey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.GetSurveyResults(userId)
}

func (s *SurveyService) CalculateAverageScore(userId int) (float64, error) {
	surveys, err := s.GetSurveyResults(userId)
	if err != nil {
		return 0, err
	}

	if len(surveys) == 0 {
		return 0, nil
	}

	var totalScore int
	for _, survey := range surveys {
		totalScore += survey.Answer
	}

	return float64(totalScore) / float64(len(surveys)), nil
}
