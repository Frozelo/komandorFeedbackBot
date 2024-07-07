package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
			b.handleMessage(update)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update)
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

	foundUser, err := b.userService.FindUser(newUser.TgId)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		msg.Text = "Произошла ошибка при проверке пользователя."
		b.api.Send(msg)
		return
	}

	if foundUser == nil {
		createdUser, err := b.userService.CreateUser(newUser)
		if err != nil {
			log.Printf("Error registering new user: %v", err)
			msg.Text = "Произошла ошибка при регистрации пользователя."
			b.api.Send(msg)
			return
		}
		log.Printf("New user registered successfully: %v", createdUser)
		msg.Text = "Привет, " + createdUser.Username + "! Ты успешно зарегистрирован. Пришли мне /start_survey, чтобы начать опрос."
	} else {
		msg.Text = "Привет, " + foundUser.Username + "! Ты уже зарегистрирован. Пришли мне /start_survey, чтобы начать опрос."
	}
	b.api.Send(msg)
}

func (b *Bot) handleStartSurvey(update tgbotapi.Update, msg tgbotapi.MessageConfig) {
	categories, err := b.surveyService.GetCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		msg.Text = "Произошла ошибка при получении категорий. Попробуйте еще раз."
		b.api.Send(msg)
		return
	}

	if len(categories) == 0 {
		msg.Text = "Нет доступных категорий для опроса."
		b.api.Send(msg)
		return
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		callbackData := fmt.Sprintf("category:%d", category.Id)
		button := tgbotapi.NewInlineKeyboardButtonData(category.Name, callbackData)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	msg.Text = "Выберите категорию для опроса:"
	msg.ReplyMarkup = replyMarkup

	b.api.Send(msg)
}

func (b *Bot) handleCallback(update tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data

	parts := strings.Split(callbackData, ":")
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат данных колбэка.")
		b.api.Send(msg)
		return
	}

	switch parts[0] {
	case "category":
		categoryID, err := strconv.Atoi(parts[1])
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат ID категории.")
			b.api.Send(msg)
			return
		}
		b.handleCategorySelection(update, categoryID)
	case "question":
		if len(parts) != 4 {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат данных колбэка.")
			b.api.Send(msg)
			return
		}
		surveyID, err := strconv.Atoi(parts[1])
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат ID опроса.")
			b.api.Send(msg)
			return
		}
		questionID, err := strconv.Atoi(parts[2])
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат ID вопроса.")
			b.api.Send(msg)
			return
		}
		answer, err := strconv.Atoi(parts[3])
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неверный формат ответа.")
			b.api.Send(msg)
			return
		}
		b.handleAnswer(update, surveyID, questionID, answer)
	}
}

func (b *Bot) handleCategorySelection(update tgbotapi.Update, categoryID int) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	questions, err := b.surveyService.GetQuestionsByCategory(categoryID)
	if err != nil {
		log.Printf("Error fetching questions: %v", err)
		msg.Text = "Произошла ошибка при получении вопросов. Попробуйте еще раз."
		b.api.Send(msg)
		return
	}

	if len(questions) == 0 {
		msg.Text = "Нет доступных вопросов для выбранной категории."
		b.api.Send(msg)
		return
	}

	// Создать новый опрос
	newSurvey := entity.Survey{
		UserId: int(update.CallbackQuery.From.ID),
	}
	createdSurvey, err := b.surveyService.CreateSurvey(newSurvey)
	if err != nil {
		log.Printf("Error creating new survey: %v", err)
		msg.Text = "Произошла ошибка при создании нового опроса. Попробуйте еще раз."
		b.api.Send(msg)
		return
	}

	// Отредактировать сообщение о выборе категории на новый текст
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Выбрана категория для опроса.")

	// Показать первый вопрос
	firstQuestion := questions[0]
	b.sendQuestion(update, createdSurvey.Id, firstQuestion)

	// Отправить запрос на редактирование сообщения
	b.api.Send(editMsg)
}

func (b *Bot) sendQuestion(update tgbotapi.Update, surveyID int, question entity.Question) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, question.Text)

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for i := 1; i <= 5; i++ {
		callbackData := fmt.Sprintf("question:%d:%d:%d", surveyID, question.Id, i)
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i), callbackData)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	msg.ReplyMarkup = replyMarkup

	b.api.Send(msg)
}

func (b *Bot) handleAnswer(update tgbotapi.Update, surveyID, questionID, answer int) {
	err := b.surveyService.SaveAnswer(entity.Answer{
		SurveyID:   surveyID,
		QuestionID: questionID,
		Answer:     answer,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка при сохранении ответа.")
		b.api.Send(msg)
		return
	}

	questions, err := b.surveyService.GetQuestionsByCategory(2)
	fmt.Printf("Questions is %v", questions)
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка при получении вопросов.")
		b.api.Send(msg)
		return
	}

	var nextQuestion *entity.Question
	for i, question := range questions {
		if question.Id == questionID && i+1 < len(questions) {
			nextQuestion = &questions[i+1]
			break
		}
	}

	if nextQuestion != nil {

		b.sendQuestion(update, surveyID, *nextQuestion)
	} else {
		avgScore, err := b.surveyService.CalculateAverageScore(surveyID)
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка при расчете среднего балла.")
			b.api.Send(msg)
			return
		}

		err = b.surveyService.SaveAvgScore(surveyID, avgScore)
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка при сохранении среднего балла.")
			b.api.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Опрос завершен! Ваш средний балл: %.2f", avgScore))
		b.api.Send(msg)
	}
}
