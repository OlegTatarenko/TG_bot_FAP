package main

import (
	"TG_bot_FAP/perm"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	kbrdMain = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.OrderStatus)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.WriteToManager),
			tgbotapi.NewKeyboardButton(perm.CallTheOffice)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(perm.GetCourier),
			tgbotapi.NewKeyboardButton(perm.RecordInService)),
	)
	/*
		kbrdYNOrg = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(perm.Organization),
				tgbotapi.NewKeyboardButton(perm.NotOrganization)),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(perm.Back)),
		)
	*/
	btnURL = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(perm.GoToChat, perm.ChatURL),
		))

	kbrdYNOrg = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Organization, perm.Organization),
			tgbotapi.NewInlineKeyboardButtonData(perm.NotOrganization, perm.NotOrganization),
		))
	numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Google", "https://www.google.ru/"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
	)
)

func main() {
	bot, err := tgbotapi.NewBotAPI(perm.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Text {
			case "/start":
				msg.Text = "Воспользуйтесь моей встроенной клавиатурой"
				msg.ReplyMarkup = kbrdMain
			case "/close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			case perm.OrderStatus:
				msg.Text = perm.InDev
			case perm.WriteToManager:
				msg.Text = "Нажмите эту кнопку для перехода в чат с менеджером"
				msg.ReplyMarkup = btnURL
			case perm.CallTheOffice:
				msg.Text = perm.CallThisNumber
			case perm.GetCourier:
				msg.Text = "Вы представляете учреждение?"
				msg.ReplyMarkup = kbrdYNOrg
			case perm.RecordInService:
				msg.Text = perm.InDev
			default:
				msg.Text = "Я тебя не понимаю"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}

/*
	for update := range updates {
		if update.Message == nil { // игнорируем отсутствие сообщения в обновлении
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		//Отображение главного меню после ввода команды /start
		switch update.Message.Text {
		case "/start":
			msg.Text = "Ты учреждение?"
			msg.ReplyMarkup = kbrdYNOrg
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		default:
			msg.Text = "Я тебя не понимаю"
		}

		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			buttonData := update.CallbackQuery.Data

			switch buttonData {
			case perm.Organization:
				msg.Text = perm.NameOfTheOrganization
			case perm.NotOrganization:
				msg.Text = perm.Adress
				//msg.ChatID = chatID
			}
		}
		// Отправляем сообщение
		if _, err = bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

/*
   		if update.CallbackQuery != nil {
   			buttonData := update.CallbackQuery.Data
   			//chatID := update.CallbackQuery.Message.Chat.ID

   			switch buttonData {
   			case perm.Organization:
   				msg.Text = perm.NameOfTheOrganization

   				//msg.ChatID = chatID
   				if _, err := bot.Send(msg); err != nil {
   					log.Panic(err)
   				}

   			case perm.NotOrganization:
   				msg.Text = perm.Adress
   				//msg.ChatID = chatID
   			}
   		}

   		if _, err := bot.Send(msg); err != nil {
   			log.Panic(err)
   		}
   	}
   }

   /*
   for update := range updates {
   		if update.Message == nil { // ignore non-Message updates
   			continue
   		}

   		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

   		//Отображение главного меню после ввода команды /start
   		switch update.Message.Text {
   		case "/start":
   			msg.Text = "Воспользуйтесь моей встроенной клавиатурой"
   			msg.ReplyMarkup = kbrdMain
   		case "/close":
   			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
   		}

   		//Статус заказа, написать менеджеру, звонок в офис
   		switch update.Message.Text {
   		case perm.OrderStatus:
   			msg.Text = perm.InDev
   		case perm.WriteToManager:
   			msg.Text = "Нажмите эту кнопку для перехода в чат с менеджером"
   			msg.ReplyMarkup = btnURL
   			//msg.Text = ""
   		case perm.CallTheOffice:
   			msg.Text = perm.CallThisNumber
   		case perm.RecordInService:
   			msg.Text = perm.InDev
   		}

   		//Обработка нажатия inlinе-кнопок
   		if update.CallbackQuery != nil {
   			buttonData := update.CallbackQuery.Data
   			chatID := update.CallbackQuery.Message.Chat.ID

   			switch buttonData {
   			case perm.Organization:
   				msg.Text = perm.NameOfTheOrganization
   				msg.ChatID = chatID
   				//if _, err := bot.Send(msg); err != nil {
   				//	log.Panic(err)
   				//}

   			case perm.NotOrganization:
   				msg.Text = perm.Adress
   				msg.ChatID = chatID
   			}
   		}

   		//Вызвать курьера
   		switch update.Message.Text {
   		case perm.GetCourier:
   			msg.Text = "Вы представляете учреждение?"
   			msg.ReplyMarkup = kbrdYNOrg

   		case perm.NotOrganization:
   			msg.Text = "Укажите полный адрес, где забрать"
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал адрес")),
   			)
   		case perm.Organization:
   			msg.Text = perm.NameOfTheOrganization
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал учреждение")),
   			)
   		case "Указал учреждение":
   			msg.Text = "Укажите полный адрес, где забрать"
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал адрес")),
   			)
   		case "Указал адрес":
   			msg.Text = "Укажите имя контактного лица"
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал имя контактного лица")),
   			)
   		case "Указал имя контактного лица":
   			msg.Text = "Нажмите кнопку, чтобы указать свой номер телефона или введите другой номер телефона"
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал номер телефона")),
   			)
   		case "Указал номер телефона":
   			msg.Text = "Для чего вызываете курьера? Например, ..."
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал цель вызова курьера")),
   			)
   		case "Указал цель вызова курьера":
   			msg.Text = "Укажите удобное время приезда курьера"
   			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   				tgbotapi.NewKeyboardButtonRow(
   					tgbotapi.NewKeyboardButton("Указал время")),
   			)
   		case "Указал время":
   			msg.Text = "Спасибо, информация записана. Для подтверждения с вами свяжется менеджер в ближайшее время."
   			msg.ReplyMarkup = kbrdMain

   		case perm.Back:
   			msg.ReplyMarkup = kbrdMain

   		}

   		/*
   			//Вызвать курьера
   			switch update.Message.Text {
   			case perm.GetCourier:
   				msg.Text = "Вы представляете учреждение?"
   				msg.ReplyMarkup = kbrdYNOrg

   			case perm.NotOrganization:
   				msg.Text = "Укажите полный адрес, где забрать"
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал адрес")),
   				)
   			case perm.Organization:
   				msg.Text = perm.NameOfTheOrganization
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал учреждение")),
   				)
   			case "Указал учреждение":
   				msg.Text = "Укажите полный адрес, где забрать"
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал адрес")),
   				)
   			case "Указал адрес":
   				msg.Text = "Укажите имя контактного лица"
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал имя контактного лица")),
   				)
   			case "Указал имя контактного лица":
   				msg.Text = "Нажмите кнопку, чтобы указать свой номер телефона или введите другой номер телефона"
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал номер телефона")),
   				)
   			case "Указал номер телефона":
   				msg.Text = "Для чего вызываете курьера? Например, ..."
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал цель вызова курьера")),
   				)
   			case "Указал цель вызова курьера":
   				msg.Text = "Укажите удобное время приезда курьера"
   				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
   					tgbotapi.NewKeyboardButtonRow(
   						tgbotapi.NewKeyboardButton("Указал время")),
   				)
   			case "Указал время":
   				msg.Text = "Спасибо, информация записана. Для подтверждения с вами свяжется менеджер в ближайшее время."
   				msg.ReplyMarkup = kbrdMain

   			case perm.Back:
   				msg.ReplyMarkup = kbrdMain

   			}

   if _, err := bot.Send(msg); err != nil {
   log.Panic(err)
   }

   }
*/
