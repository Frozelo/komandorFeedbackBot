package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BotHandler struct {
	bot *tgbotapi.BotAPI
}

func NewBotHandler(bot *tgbotapi.BotAPI) *BotHandler {
	return &BotHandler{bot: bot}
}

func (bh *BotHandler) HandleCommands(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "You just entered the start command!"))
	case "test":
		bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Wow! Very rare command!"))

	default:
		bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Unknown command!"))
	}
}
