package main

import (
	"context"
	"log"

	"github.com/Frozelo/komandorFeedbackBot/internal/bot"
	postgres_storage "github.com/Frozelo/komandorFeedbackBot/internal/storage/postgres"
)

func main() {
	cfgPath := "internal/bot/config.yml"
	cfg, err := bot.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bot configuration loaded from %s", cfgPath)

	db, err := postgres_storage.NewDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting postgres")

	defer db.Close(context.Background())

	b, err := bot.NewBot(db, cfg.Bot.ApiKey)

	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
