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

	var kbrdMain = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏃Вызвать курьера"),
			tgbotapi.NewKeyboardButton("❓Статус заказа")),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏Написать менеджеру"),
			tgbotapi.NewKeyboardButton("📞Позвонить в офис"),
			tgbotapi.NewKeyboardButton("🔧Записаться в Сервис")),
	)

	var kbrdYN = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Я представляю учреждение🏢"),
			tgbotapi.NewKeyboardButton("Нет, я не учреждение🙂")),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад🔙")),
	)

	var btnURL = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("➡Нажмите для перехода в чат⬅", "https://t.me/ultrafilaret"),
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
		case "❓Статус заказа":
			msg.Text = "Извините, эта функция еще в разработке🤷‍♂️"
		case "✏Написать менеджеру":
			msg.ReplyMarkup = btnURL
			//msg.Text = "Перейти в чат с оператором: @ultrafilaret  (https://t.me/ultrafilaret)"
		case "📞Позвонить в офис":
			msg.Text = "Для связи с офисом позвоните, пожалуйста, по этому номеру +79045340560"
		case "🔧Записаться в Сервис":
			msg.Text = "Извините, эта функция еще в разработке🤷‍♂️"
		}

		switch update.Message.Text {
		case "🏃Вызвать курьера":
			msg.ReplyMarkup = kbrdYN
		case "Я представляю учреждение🏢":
			msg.Text = "Укажите наименование учреждения"
		case "Нет, я не учреждение🙂":
			msg.Text = "Извините, эта функция еще в разработке🤷‍♂️"
		case "Назад🔙":
			msg.ReplyMarkup = kbrdMain
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
