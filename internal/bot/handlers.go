package bot

import (
	"log"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api           *tgbotapi.BotAPI
	surveyService *service.SurveyService
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = true

	return &Bot{
		api:           api,
		surveyService: service.NewSurveyService(),
	}, nil
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	chatID := update.Message.Chat.ID

	switch update.Message.Text {
	case "/start":
		b.surveyService.StartSurvey(int(chatID))
		nextQuestion := b.surveyService.GetNextQuestion(int(chatID))
		if nextQuestion != nil {
			msg.Text = nextQuestion.Text
		} else {
			msg.Text = "Нет доступных вопросов."
		}
	default:
		b.surveyService.AnswerQuestion(int(chatID), update.Message.Text)
		nextQuestion := b.surveyService.GetNextQuestion(int(chatID))
		if nextQuestion != nil {
			msg.Text = nextQuestion.Text
		} else {
			msg.Text = "Спасибо за ваши ответы!"
		}
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update)
		}
	}
}
