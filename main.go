package main

import (
	"TG_bot_FAP/perm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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

// Forms хеш-таблица для хранения данных пользователей
var Forms = make(map[int64]UserData)

// UserData структура для записи данных пользователя при вызове курьера или при записи на сервис
type UserData struct {
	UserID           int64
	DataGetCourier   Blank
	DataRecInService string
	Index            int
}

// Blank структура для записи данных пользователя при вызове курьера
type Blank struct {
	Intro   string
	Org     string
	Address string
	Person  string
	Phone   string
	Purpose string
	Time    string
}

var UD = UserData{
	0,
	B,
	"",
	0,
}

var B = Blank{
	"Спасибо! Вы ввели следующие данные: ",
	"Учреждение (если применимо): ",
	"Адрес, где забрать: ",
	"Контактное лицо: ",
	"Контактный телефон: ",
	"Цель вызова курьера: ",
	"Дата, время приезда курьера: ",
}

// WriteUDStart функция для создания записи ползователя в мапе с ID ползователя
func WriteUDStart(ID int64) {
	if _, ok := Forms[ID]; !ok {
		//если нет ключа, равного ID, то создаем элемент мапы, записав в соотв. поле ID
		UD.UserID = ID
		Forms[ID] = UD
	}
}

// WriteUDIndex функция для изменения индекса в данных пользователя
func WriteUDIndex(ID int64) {
	//переменная с копией структуры, чтобы обратиться к полю структуры внутри мапы, т.к. обращение через Forms[ID].index = i языком не предусмотрено
	//изменяем счетчик и записываем его в index
	temp := Forms[ID]
	temp.Index = i
	Forms[ID] = temp
}

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
	//ClientForm := perm.Form

	// Loop through each update.
	for update := range updates {

		// Создаем запись о пользователе в мапе, если записи не существует
		ID := update.Message.Chat.ID
		Text := update.Message.Text
		WriteUDStart(ID)
		//переменная с копией структуры, чтобы получить значение поля структуры из мапы, т.к. обращение через Forms[ID].index > 0 языком не предусмотрено
		ind := Forms[ID].Index

		// Опрос клиента после нажатия кнопки Вызвать курьера, т.е. при условии, что i > 0
		if ind > 0 && update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch ind {
			case 1:
				//записываем полученное наименование организации в поле Org
				temp := Forms[ID]
				temp.DataGetCourier.Org = temp.DataGetCourier.Org + Text
				Forms[ID] = temp
				i++ //2
				WriteUDIndex(ID)
				msg.Text = perm.Adress
			case 2:
				//записываем полученный адрес в поле Address
				temp := Forms[ID]
				temp.DataGetCourier.Address = temp.DataGetCourier.Address + Text
				Forms[ID] = temp
				i++ //3
				WriteUDIndex(ID)
				msg.Text = "Укажите имя контактного лица"
			case 3:
				//записываем полученное имя в поле Person
				temp := Forms[ID]
				temp.DataGetCourier.Person = temp.DataGetCourier.Person + Text
				Forms[ID] = temp
				i++ //4
				WriteUDIndex(ID)
				msg.ReplyMarkup = btnContact
				msg.Text = "Введите номер телефона или нажмите кнопку, чтобы отправить номер телефона, к которому привязан ваш аккаунт телеграма"

			case 4:
				//записываем полученный телефон в поле Phone
				temp := Forms[ID]
				temp.DataGetCourier.Phone = temp.DataGetCourier.Phone + Text
				Forms[ID] = temp
				i++ //5
				WriteUDIndex(ID)
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Для чего вызываете курьера? Например, часто указывают:" +
					"\n- забрать картриджи, заправить, вернуть." +
					"\n- купить 1 новый картридж ce285a и выставить счет." +
					"\n- после заправки привезти акт сверки и новый счет." +
					"\n- забрать подписанный договор." +
					"\n...или ваш вариант"
			case 5:
				//записываем полученную цель вызова в поле Purpose
				temp := Forms[ID]
				temp.DataGetCourier.Purpose = temp.DataGetCourier.Purpose + Text
				Forms[ID] = temp
				i++ //6
				WriteUDIndex(ID)
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Укажите удобную дату и время приезда курьера"
			case 6:
				//записываем полученную дату и время в поле Time
				temp := Forms[ID]
				temp.DataGetCourier.Time = temp.DataGetCourier.Time + Text
				Forms[ID] = temp
				i++ //7
				WriteUDIndex(ID)
				msg.ReplyMarkup = kbrdYN
				//Выводим записанные данные клиента в виде сообщения
				//TODO: заменить temp.DataGetCourier.Time на вывод всех данных структуры DataGetCourier, сейчас выводит только время
				msg.Text = temp.DataGetCourier.Time + "\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
			case 7:
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки.\n "
				i = 0
				WriteUDIndex(ID)
			}
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.Message != nil {
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
			default:
				msg.Text = "Я тебя не понимаю"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		}
		if update.CallbackQuery != nil {
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
				ID := update.CallbackQuery.Message.Chat.ID
				WriteUDIndex(ID)
			case perm.NotOrganization:
				msg.Text = perm.Adress
				i = 2
				ID := update.CallbackQuery.Message.Chat.ID
				WriteUDIndex(ID)
			case perm.Yes:
				//ClientForm = perm.Form - обнуление, которое не работает
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки."
				i = 0
				ID := update.CallbackQuery.Message.Chat.ID
				WriteUDIndex(ID)
			case perm.No:
				//ClientForm = perm.Form - обнуление, которое не работает
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
				i = 0
				ID := update.CallbackQuery.Message.Chat.ID
				WriteUDIndex(ID)
			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)

			}
		}

	}
}
