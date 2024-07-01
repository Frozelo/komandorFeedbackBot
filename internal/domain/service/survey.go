package service

import (
	"sync"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
)

type SurveyService struct {
	mu      sync.Mutex
	surveys map[int]*entity.Survey
}

func NewSurveyService() *SurveyService {
	return &SurveyService{surveys: make(map[int]*entity.Survey)}
}

func (ss *SurveyService) StartSurvey(chatId int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.surveys[chatId] = &entity.Survey{
		Questions: []entity.Question{
			{Text: "Как вам наш продукт?", Answered: false},
			{Text: "Какие функции вы бы хотели видеть?", Answered: false},
		},
	}
}

func (ss *SurveyService) GetNextQuestion(chatId int) *entity.Question {
	ss.mu.Lock()

	defer ss.mu.Unlock()

	survey, exists := ss.surveys[chatId]

	if !exists {
		return nil
	}

	for i, q := range survey.Questions {
		if !q.Answered {
			return &survey.Questions[i]
		}
	}

	return nil
}

func (s *SurveyService) AnswerQuestion(chatID int, answer string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	survey, exists := s.surveys[chatID]
	if !exists {
		return
	}

	for i, q := range survey.Questions {
		if !q.Answered {
			survey.Questions[i].Answered = true
			survey.Questions[i].Answer = answer
			break
		}
	}
}
