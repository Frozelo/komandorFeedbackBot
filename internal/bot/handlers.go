package bot

import (
	"log"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
	"github.com/Frozelo/komandorFeedbackBot/internal/domain/service"
	"github.com/Frozelo/komandorFeedbackBot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

type BotConfig struct {
	ApiKey string `yaml:"api_key"`
}

type Bot struct {
	api         *tgbotapi.BotAPI
	userService *service.UserService
}

func NewBot(db *pgx.Conn, apiKey string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	return &Bot{
		api:         bot,
		userService: userService,
	}, nil
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

func (b *Bot) handleMessage(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Text {
	case "/start":
		b.handleStart(update, msg)
	case "/start_survey":
		b.handleStartSurvey(update, msg)
	case "/help":
		b.handleHelp(update, msg)
	default:
		msg.Text = "Неизвестная команда. Напишите /help для списка доступных команд."
		b.api.Send(msg)
	}
}

func (b *Bot) handleStart(update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	newUser := entity.User{
		TgId:     int(update.Message.From.ID),
		Username: update.Message.From.UserName,
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

func (b *Bot) handleStartSurvey(update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	msg.Text = "Опрос начат. Пожалуйста, ответьте на следующие вопросы..."
	b.api.Send(msg)
	panic("implement me")
}

func (b *Bot) handleHelp(update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	msg.Text = "Список доступных команд:\n/start - Регистрация пользователя\n/start_survey - Начать опрос\n/help - Список команд"
	b.api.Send(msg)
	panic("implement me")
}
