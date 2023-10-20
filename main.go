package main

import (
	"TG_bot_FAP/perm"
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

	var kbrdMain = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.GetCourier),
			tgbotapi.NewKeyboardButton(perm.OrderStatus)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.WriteToManager),
			tgbotapi.NewKeyboardButton(perm.CallTheOffice),
			tgbotapi.NewKeyboardButton(perm.RecordInService)),
	)

	var kbrdYNOrg = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.Organization),
			tgbotapi.NewKeyboardButton(perm.NotOrganization)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.Back)),
	)

	var btnURL = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(perm.GoToChat, perm.ChatURL),
		),
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
			msg.ReplyMarkup = kbrdMain
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		}

		switch update.Message.Text {
		case perm.OrderStatus:
			msg.Text = perm.InDev
		case perm.WriteToManager:
			msg.ReplyMarkup = btnURL
		case perm.CallTheOffice:
			msg.Text = perm.CallThisNumber
		case perm.RecordInService:
			msg.Text = perm.InDev
		}

		switch update.Message.Text {
		case perm.GetCourier:
			msg.ReplyMarkup = kbrdYNOrg
		case perm.Organization:
			msg.Text = perm.NameOfTheOrganization
		case perm.NotOrganization:
			msg.Text = perm.InDev
		case perm.Back:
			msg.ReplyMarkup = kbrdMain
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
