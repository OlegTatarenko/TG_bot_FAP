package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6459859641:AAEyvBu87_bFrIjsJumyv4KA0gpF1GR6d94")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вызвать курьера"),
			tgbotapi.NewKeyboardButton("Статус заказа")),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg.ReplyMarkup = keyboard
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
