package main

import (
	"log"

	"github.com/Frozelo/komandorFeedbackBot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfgPath := "internal/bot/config.yml"
	cfg, err := bot.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.Bot.ApiKey)
	if err != nil {
		log.Panic(err)
	}

	botApi.Debug = true

	log.Printf("Authorized on account %s", botApi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botApi.GetUpdatesChan(u)

	botHandler := bot.NewBotHandler(botApi)
	for update := range updates {
		if update.Message != nil {
			botHandler.HandleCommands(update.Message)
		}
	}
}
