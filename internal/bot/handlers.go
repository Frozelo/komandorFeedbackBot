package bot

import (
	"fmt"
	"log"
	"strconv"

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
	api           *tgbotapi.BotAPI
	userService   *service.UserService
	surveyService *service.SurveyService
}

func NewBot(db *pgx.Conn, apiKey string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	surveyRepo := repository.NewSurveyRepository(db)
	surveyService := service.NewSurveyService(surveyRepo)

	return &Bot{
		api:           bot,
		userService:   userService,
		surveyService: surveyService,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleCommand(update)
		}
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update)
		}
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		b.handleStart(update)
	case "survey":
		b.handleSurvey(update)
	default:
		msg.Text = "Неизвестная команда. Напишите /help для списка доступных команд."
		b.api.Send(msg)
	}
}

func (b *Bot) handleStart(update tgbotapi.Update) {
	log.Printf("handleStart called with update: %v", update)

	newUser := entity.User{
		TgId:     int(update.Message.From.ID),
		Username: update.Message.From.UserName,
	}

	log.Printf("Constructed newUser: %v", newUser)

	existingUser, err := b.findUserByTgId(newUser.TgId)
	log.Printf("Existing user check completed, existingUser: %v, err: %v", existingUser, err)
	if err != nil {
		log.Printf("Error finding user by TgId: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка. Пожалуйста, попробуйте еще раз.")
		b.api.Send(msg)
		return
	}

	if existingUser != nil {
		log.Printf("Existing user found: %v", existingUser)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже зарегистрированы. Ваш никнейм: "+existingUser.Username)
		b.api.Send(msg)
		return
	}

	log.Printf("No existing user found, proceeding to register new user")

	createdUser, err := b.registerNewUser(newUser)
	log.Printf("User registration completed, createdUser: %v, err: %v", createdUser, err)
	if err != nil {
		log.Printf("Error registering new user: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при регистрации. Пожалуйста, попробуйте еще раз.")
		b.api.Send(msg)
		return
	}

	log.Printf("New user registered successfully: %v", createdUser)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+createdUser.Username+"! Ты успешно зарегистрирован. Пришли мне /start_survey, чтобы начать опрос.")
	b.api.Send(msg)
}

func (b *Bot) handleSurvey(update tgbotapi.Update) {
	question := "Как вы оцениваете качество нашего сервиса от 0 до 5?"
	callbackButtons := make([]tgbotapi.InlineKeyboardButton, 6)
	for i := 0; i <= 5; i++ {
		callbackButtons[i] = tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), strconv.Itoa(i))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, question)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(callbackButtons...))
	b.api.Send(msg)
}

func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	rating, err := strconv.Atoi(callbackQuery.Data)
	if err != nil {
		log.Printf("Invalid callback data: %v", err)
		return
	}

	survey := entity.Survey{
		UserId:   int(callbackQuery.From.ID),
		Question: "Как вы оцениваете качество нашего сервиса от 0 до 5?",
		Answer:   rating,
	}

	_, err = b.surveyService.CreateSurvey(survey)
	if err != nil {
		log.Printf("Error creating survey: %v", err)
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Ошибка при сохранении ответа. Попробуйте еще раз.")
		b.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Ваш ответ сохранен. Спасибо за участие!")

	averageScore, err := b.surveyService.CalculateAverageScore(survey.UserId)
	if err != nil {
		log.Printf("Error calculating average score: %v", err)
		msg.Text = "Ошибка при вычислении среднего результата. Попробуйте еще раз."
		b.api.Send(msg)
		return
	}

	msg.Text = fmt.Sprintf("Средняя оценка вашего опроса: %.2f", averageScore)
	b.api.Send(msg)
}

func (b *Bot) findUserByTgId(tgId int) (*entity.User, error) {
	log.Printf("findUserByTgId called with tgId: %d", tgId)
	existingUser, err := b.userService.FindUser(tgId)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}

	log.Printf("User found: %v", existingUser)
	return existingUser, nil
}

func (b *Bot) registerNewUser(newUser entity.User) (*entity.User, error) {
	log.Printf("registerNewUser called with newUser: %v", newUser)
	createdUser, err := b.userService.CreateUser(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Printf("User created: %v", createdUser)
	return createdUser, nil
}
