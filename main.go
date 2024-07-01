package main

import (
	"log"

	"github.com/Frozelo/komandorFeedbackBot/internal/bot"
)

func main() {
	cfgPath := "internal/bot/config.yml"
	cfg, err := bot.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.NewBot(cfg.Bot.ApiKey)

	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
