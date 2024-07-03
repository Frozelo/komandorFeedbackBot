package bot

import (
	"log"
	"time"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
	"github.com/Frozelo/komandorFeedbackBot/internal/domain/service"
	"github.com/Frozelo/komandorFeedbackBot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

type Bot struct {
	api           *tgbotapi.BotAPI
	userService   *service.UserService
	surveyService *service.SurveyService
}

func NewBot(db *pgx.Conn, token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)

	api.Debug = true

	return &Bot{
		api:           api,
		userService:   service.NewUserService(userRepo),
		surveyService: service.NewSurveyService(),
	}, nil
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Text {
	case "/start":
		joinedAt := time.Unix(int64(update.Message.Date), 0).Format("2006-01-02")
		newUser := entity.User{
			TgId:     int(update.Message.From.ID),
			Username: update.Message.From.UserName,
			JoinedAt: joinedAt,
		}
		user, err := b.userService.CreateUser(newUser)

		if err != nil {
			msg.Text = "Ошибка при создании пользователя."
			log.Printf("Error creating user: %v", err)

		} else {
			msg.Text = "Привет, " + user.Username + "! Ты уже зарегистрирован. Пришли мне /start_survey, чтобы начать опрос."
		}
		b.api.Send(msg)

	}
}

func (b *Bot) Start() {

	log.Printf("Bot is running...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update)
		}
	}
}
