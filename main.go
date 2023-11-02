package main

import (
	"TG_bot_FAP/perm"
	"log"
	"strings"

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

	btnContact = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(perm.Contact)))

	btnURL = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(perm.GoToChat, perm.ChatURL),
		))
	kbrdYNOrg = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Organization, perm.Organization),
			tgbotapi.NewInlineKeyboardButtonData(perm.NotOrganization, perm.NotOrganization),
		))
	kbrdYN = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Yes, perm.Yes),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, perm.No),
		))
)

// Счетчик для вывода инлайн-кнопок и записи данных пользователя
var i = 0

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

	// Срез для записи информации клиента при вызове курьера
	ClientForm := perm.Form

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// Реакция на нажатия кнопок главного меню
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
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
			case perm.RecordInService:
				msg.Text = perm.InDev
				//default:
				//	msg.Text = "Я тебя не понимаю"
			}

			// Опрос клиента после нажатия кнопки Вызвать курьера
			switch i {
			case 1:
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = perm.Adress
				i++ //2
			case 2:
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = "Укажите имя контактного лица"
				i++ //3
			case 3:
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = "Введите номер телефона или нажмите кнопку, чтобы отправить номер телефона, к которому привязан ваш аккаунт телеграма"
				msg.ReplyMarkup = btnContact
				i++ //4
			case 4:
				msg.ReplyMarkup = kbrdMain
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = "Для чего вызываете курьера? Например, часто указывают:" +
					"\n1) Забрать картриджи, заправить, вернуть." +
					"\n2) Купить 1 новый картридж ce285a и выставить счет." +
					"\n3) После заправки привезти акт сверки и новый счет." +
					"\n4) Забрать подписанный договор." +
					"\n...или ваш вариант"
				i++ //5
			case 5:
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = "Укажите удобную дату и время приезда курьера"
				i++ //6
			case 6:
				ClientForm[i] = ClientForm[i] + msg.Text
				msg.Text = strings.Join(ClientForm, "\n") + "\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
				msg.ReplyMarkup = kbrdYN
				i++ //7
			case 7:
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки.\n "
				i = 0
				ClientForm = perm.Form
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

			switch update.CallbackQuery.Data {
			case perm.Organization:
				msg.Text = perm.NameOfTheOrganization
				i = 1
			case perm.NotOrganization:
				msg.Text = perm.Adress
				i = 2
			case perm.Yes:
				ClientForm = perm.Form
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки."
				i = 0
			case perm.No:
				ClientForm = perm.Form
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
				i = 0
			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)

			}
		}
		if i == 0 {
			ClientForm = perm.Form
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
