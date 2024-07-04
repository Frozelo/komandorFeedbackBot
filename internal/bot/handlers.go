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
	default:
		msg.Text = "Неизвестная команда. Напишите /help для списка доступных команд."
		b.api.Send(msg)
	}
}

func (b *Bot) handleStart(update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	log.Printf("handleStart called with update: %v", update)

	newUser := entity.User{
		TgId:     int(update.Message.From.ID),
		Username: update.Message.From.UserName,
	}

	log.Printf("Constructed newUser: %v", newUser)

	existingUser, err := b.findUserByTgId(newUser.TgId, &msg)
	log.Printf("Existing user check completed, existingUser: %v, err: %v", existingUser, err)
	if err != nil {
		log.Printf("Error finding user by TgId: %v", err)
		b.api.Send(msg)
		return
	}

	if existingUser != nil {
		log.Printf("Existing user found: %v", existingUser)
		msg.Text = "Вы уже зарегистрированы. Ваш никнейм: " + existingUser.Username
		b.api.Send(msg)
		return
	}

	log.Printf("No existing user found, proceeding to register new user")

	createdUser, err := b.registerNewUser(newUser, &msg)
	log.Printf("User registration completed, createdUser: %v, err: %v", createdUser, err)
	if err != nil {
		log.Printf("Error registering new user: %v", err)
		b.api.Send(msg)
		return
	}

	log.Printf("New user registered successfully: %v", createdUser)

	msg.Text = "Привет, " + createdUser.Username + "! Ты успешно зарегистрирован. Пришли мне /start_survey, чтобы начать опрос."
	b.api.Send(msg)
}

func (b *Bot) findUserByTgId(tgId int, msg *tgbotapi.MessageConfig) (*entity.User, error) {
	log.Printf("findUserByTgId called with tgId: %d", tgId)
	existingUser, err := b.userService.FindUser(tgId)
	if err != nil {
		msg.Text = "Произошла ошибка при проверке существующего пользователя."
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}

	if existingUser != nil {
		log.Printf("User found: %v", existingUser)
		msg.Text = "Вы уже зарегистрированы. Ваш никнейм: " + existingUser.Username
		return existingUser, nil
	}

	log.Printf("User not found for tgId: %d", tgId)
	return nil, nil
}

func (b *Bot) registerNewUser(newUser entity.User, msg *tgbotapi.MessageConfig) (*entity.User, error) {
	log.Printf("registerNewUser called with newUser: %v", newUser)
	createdUser, err := b.userService.CreateUser(newUser)
	if err != nil {
		msg.Text = "Произошла ошибка при регистрации. Попробуйте еще раз."
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Printf("User created: %v", createdUser)
	return createdUser, nil
}
