package main

import (
	"TG_bot_FAP/perm"
	"TG_bot_FAP/remonline"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
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
			tgbotapi.NewKeyboardButton(perm.RecInService)),
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
	kbrdYNRecInService = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Yes, "yes"),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, "no"),
		))
	kbrdBeforeAfterLunch = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("До обеда (9:00-13:00)", "beforeLunch"),
			tgbotapi.NewInlineKeyboardButtonData("После обеда (13:00 - 18:00)", "afterLunch"),
		))
)

// Счетчик для вывода инлайн-кнопок и записи данных пользователя
var i = 0

// users мапа для хранения данных пользователей
var users = make(map[int64]user)

// sliceBtnString срез для хранения значений "text" кнопок в клавиатуре kbrdDate
var btnsKbrdDates = make([]string, 6, 6)

// user структура для записи данных пользователя при вызове курьера или при записи на сервис
type user struct {
	FirstName    string
	Username     string
	GetCourier   blankGetCourier
	RecInService blankRecInService
	Index        int
}

// blankGetCourier структура для записи данных пользователя при вызове курьера
type blankGetCourier struct {
	Intro   string
	Org     string
	Address string
	Person  string
	Phone   string
	Purpose string
	Time    string
}

// blankRecInService структура для записи данных пользователя при записи в сервис
type blankRecInService struct {
	Intro   string
	Date    string
	Time    string
	Device  string
	Problem string
	Phone   string
	Name    string
}

var getCourierData = blankGetCourier{
	"Спасибо! Вы ввели следующие данные: ",
	"Учреждение (если применимо): ",
	"Адрес, где забрать: ",
	"Контактное лицо: ",
	"Номер телефона: ",
	"Цель вызова курьера: ",
	"Дата, время приезда курьера: ",
}
var recInServiceData = blankRecInService{
	"Спасибо! Вы ввели следующие данные: ",
	"Дата: ",
	"Время: ",
	"Устройство: ",
	"Проблема: ",
	"Номер телефона: ",
	"Контактное лицо: ",
}

var userData = user{
	"Name: ",
	"Username: ",
	getCourierData,
	recInServiceData,
	0,
}

var t = time.Now()

// btnsMainMenu слайс с названиями кнопок главного меню
var btnsMainMenu = []string{
	perm.OrderStatus,
	perm.WriteToManager,
	perm.CallTheOffice,
	perm.GetCourier,
	perm.RecInService,
}

// writeUserData функция для создания записи пользователя в мапе с ID пользователя
func writeUserData(ID int64, FirstName, Username string) {
	if _, ok := users[ID]; !ok {
		users[ID] = userData //Если нет ключа, равного ID, то создаем элемент мапы с ключом = ID
		temp := users[ID]    //переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
		//т.к. изменение значения поля через users[ID].index = i языком не предусмотрено
		temp.FirstName = FirstName //Записываем в поле имя пользователя
		temp.Username = Username   //Записываем в поле Username пользователя
		temp.Index = i             // Меняем индекс у пользователя
		users[ID] = temp           //Изменяем запись в мапе
	} else {
		// меняем индекс у пользователя
		temp := users[ID]
		temp.Index = i
		users[ID] = temp
	}
}

// writeUserDataIndex функция для изменения индекса в данных пользователя
func writeUserDataIndex(ID int64) {
	//переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
	//т.к. изменение значения поля через users[ID].index = i языком не предусмотрено
	//счетчик записываем в index
	temp := users[ID]
	temp.Index = i
	users[ID] = temp
}

//TODO: Тест - запись в субботу в между 14:00 и 14:30

// kbrdDate формирует клавиатуру из 6 инлайн кнопок, для чего создает слайс длиной 6 из строк,
// где каждая строка - дата для инлайн кнопки на запись в сервис
func kbrdDate(t time.Time) tgbotapi.InlineKeyboardMarkup {
	j := 0                 // индексы сдвига даты
	hh, mm, _ := t.Clock() //берем часы и минуты из времени нажатия на кнопку, т.е. из текущего времени
	timeNOW := hh*60 + mm  //текущее время в минутах с начала суток
	for k := 0; k < len(btnsKbrdDates); k++ {
		l := 0 // индексы сдвига даты
		//если время нажатия на кнопку > 16:30 (cntrlTimeWorkday) или день нажатия на кнопку - суббота и время > 14:00 (cntrlTimeSat), то
		//увеличиваем l на единицу, чтобы пропустить сегодняшнюю дату
		if timeNOW > perm.CntrlTimeWorkday || (t.Weekday() == time.Saturday && timeNOW > perm.CntrlTimeSat) {
			l = 1
		}

		datesRecInService := t.Add(time.Hour * 24 * time.Duration(k+j+l)) //добавляем к текущей дате 24*(k+j+l) часов
		weekDay := datesRecInService.Weekday()                            //определяем для даты день недели

		//если день недели - воскресенье, т.е. выходной, то j=1 , чтобы пропустить дату воскресенья
		if weekDay == time.Sunday {
			j = 1
			datesRecInService = t.Add(time.Hour * 24 * time.Duration(k+j+l))
		}
		btnsKbrdDates[k] = datesRecInService.Format("02.01.2006") //форматируем дату в строку и записываем в слайс
	}

	kbrdDates := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[0], btnsKbrdDates[0]),
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[1], btnsKbrdDates[1]),
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[2], btnsKbrdDates[2]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[3], btnsKbrdDates[3]),
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[4], btnsKbrdDates[4]),
			tgbotapi.NewInlineKeyboardButtonData(btnsKbrdDates[5], btnsKbrdDates[5]),
		))
	return kbrdDates
}

// mustTokenTg - функция для получения токена телеграм бота через флаг -tgbot-token
// для запуска из командной строки необходимо:  go run TG_bot_FAP -tgbot-token 'значение токена'   или же
//   - go build (собираем exe-файл, если его еще нет, если есть - пропускаем эту команду)
//   - ./TG_bot_FAP -tgbot-token 'значение токена' (запускаем exe-файл с флагом '-tgbot-token', указывая значение токена)
func mustTokenTg() string {
	token := flag.String(
		"tgbot-token",
		"",
		"токен для доступа к телеграм боту / token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("токен не указан / token is not specified")
	}

	return *token
}

// isCommand функция, определяющая сообщения, полученные после нажатия на кнопки главного меню
func isCommand(text string, btnsMainMenu []string) bool {
	b := false
	for _, v := range btnsMainMenu {
		if text == v {
			b = true
		}
	}
	return b
}

func main() {
	bot, err := tgbotapi.NewBotAPI(mustTokenTg())
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

		//TODO удалить данные юзера из мапы и отрисовать ему главное меню, если он не совершает действий какое-то время

		//TODO дописать проверку корректности ввода номера телефона при вызове курьера и записи в сервис, взять готовую функцию

		// Опрос клиента после нажатия кнопок Вызвать курьера, Запись в сервис или Статус заказа, т.е. при условии, что user.Index > 0
		if update.Message != nil && users[update.Message.From.ID].Index > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			ID := update.Message.Chat.ID
			Text := update.Message.Text

			switch users[update.Message.From.ID].Index {
			case 1: //ветка после нажатия на кнопку Вызвать курьера
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForGetCourier
				} else {
					//записываем полученное наименование организации в поле Org
					temp := users[ID]
					temp.GetCourier.Org = temp.GetCourier.Org + Text
					users[ID] = temp
					i = 2
					writeUserDataIndex(ID)
					msg.Text = perm.Address
				}
			case 2:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForGetCourier
				} else {
					//записываем полученный адрес в поле Address
					temp := users[ID]
					temp.GetCourier.Address = temp.GetCourier.Address + Text
					users[ID] = temp
					i = 3
					writeUserDataIndex(ID)
					msg.Text = "Укажите имя контактного лица"
				}
			case 3:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForGetCourier
				} else {
					//записываем полученное имя в поле Person
					temp := users[ID]
					temp.GetCourier.Person = temp.GetCourier.Person + Text
					users[ID] = temp
					i = 4
					writeUserDataIndex(ID)
					msg.ReplyMarkup = btnContact
					msg.Text = perm.EnterPhone
				}
			case 4:
				//определяем, отправлен телефон текстовым сообщением или пользователь нажал инлайн кнопку "Отправить телефон"
				phone := "_"
				if update.Message.Text != "" {
					phone = update.Message.Text
				} else if update.Message.Contact != nil {
					phone = update.Message.Contact.PhoneNumber
				}
				//проверяем корректность введенного номера телефона
				if !remonline.IsValidPhone(phone) {
					msg.Text = perm.NotCorrectPhone
					msg.ReplyMarkup = btnContact
				} else {
					//записываем полученный телефон в поле Phone
					temp := users[ID]
					if update.Message.Text != "" {
						//если пользователь отправил телефон текстовым сообщением, то записываем его в поле
						temp.GetCourier.Phone = temp.GetCourier.Phone + Text
					} else {
						//если пользователь нажал инлайн кнопку "Отправить телефон", то записываем его в поле
						temp.GetCourier.Phone = temp.GetCourier.Phone + update.Message.Contact.PhoneNumber
					}
					users[ID] = temp
					i = 5
					writeUserDataIndex(ID)
					msg.ReplyMarkup = kbrdMain
					msg.Text = "Для чего вызываете курьера? Например, часто указывают:" +
						"\n- забрать картриджи, заправить, вернуть." +
						"\n- купить 1 новый картридж ce285a и выставить счет." +
						"\n- после заправки привезти акт сверки и новый счет." +
						"\n- забрать подписанный договор." +
						"\n...или ваш вариант"
				}
			case 5:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForGetCourier
				} else {
					//записываем полученную цель вызова в поле Purpose
					temp := users[ID]
					temp.GetCourier.Purpose = temp.GetCourier.Purpose + Text
					users[ID] = temp
					i = 6
					writeUserDataIndex(ID)
					msg.ReplyMarkup = kbrdMain
					msg.Text = "Укажите удобную дату и время приезда курьера"
				}
			case 6:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForGetCourier
				} else {
					//записываем полученную дату и время в поле Time
					temp := users[ID]
					temp.GetCourier.Time = temp.GetCourier.Time + Text
					users[ID] = temp
					i = 7
					writeUserDataIndex(ID)
					msg.ReplyMarkup = kbrdYN
					//Выводим записанные данные клиента в виде сообщения
					msg.Text = users[ID].GetCourier.Intro +
						"\n" + users[ID].GetCourier.Org +
						"\n" + users[ID].GetCourier.Address +
						"\n" + users[ID].GetCourier.Person +
						"\n" + users[ID].GetCourier.Phone +
						"\n" + users[ID].GetCourier.Purpose +
						"\n" + users[ID].GetCourier.Time +
						"\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
				}
			case 7:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = "Подтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись, заполняя заявку на вызов курьера"
					msg.ReplyMarkup = kbrdYN
				} else {
					//если НЕ НАЖАТА кнопка "Да, все верно" и отправлено какое-либо сообщение, то отправляем сообщение менеджеру с записанными данными
					msg.ChatID = perm.ManagerID
					msg.Text = "⚡ Заявка на вызов курьера:" +
						"\n\n" + users[ID].FirstName +
						"\n" + users[ID].Username +
						"\n" + users[ID].GetCourier.Org +
						"\n" + users[ID].GetCourier.Address +
						"\n" + users[ID].GetCourier.Person +
						"\n" + users[ID].GetCourier.Phone +
						"\n" + users[ID].GetCourier.Purpose +
						"\n" + users[ID].GetCourier.Time +
						"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
					if _, err := bot.Send(msg); err != nil {
						panic(err)

					}
					//удаляем из мапы запись с данными пользователя
					delete(users, ID)
					//отправляем сообщение пользователю
					msg.ChatID = ID
					msg.ReplyMarkup = kbrdMain
					msg.Text = perm.Ok
				}
			case 8: //ветка после нажатия на кнопку Запись в сервис
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForRecInService
				} else {
					//записываем проблему в поле Problem
					temp := users[ID]
					temp.RecInService.Problem = temp.RecInService.Problem + Text
					users[ID] = temp
					i = 9
					writeUserDataIndex(ID)
					msg.Text = "Какое у вас устройство? Например, ноутбук ACER модель ABC."
				}
			case 9:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForRecInService
				} else {
					//записываем устройство в поле Device
					temp := users[ID]
					temp.RecInService.Device = temp.RecInService.Device + Text
					users[ID] = temp
					i = 10
					writeUserDataIndex(ID)
					msg.Text = "Укажите ваше имя"
				}
			case 10:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = perm.NotAllAnswersForRecInService
				} else {
					//записываем имя в поле Name
					temp := users[ID]
					temp.RecInService.Name = temp.RecInService.Name + Text
					users[ID] = temp
					i = 11
					writeUserDataIndex(ID)
					msg.ReplyMarkup = btnContact
					msg.Text = perm.EnterPhone
				}
			case 11:
				//определяем, отправлен телефон текстовым сообщением или пользователь нажал инлайн кнопку "Отправить телефон"
				phone := "_"
				if update.Message.Text != "" {
					phone = update.Message.Text
				} else if update.Message.Contact != nil {
					phone = update.Message.Contact.PhoneNumber
				}
				//проверяем корректность введенного номера телефона
				if !remonline.IsValidPhone(phone) {
					msg.Text = perm.NotCorrectPhone
					msg.ReplyMarkup = btnContact
				} else {
					//записываем полученный телефон в поле Phone
					temp := users[ID]
					if update.Message.Text != "" {
						//если пользователь отправил телефон текстовым сообщением, то записываем его в поле
						temp.RecInService.Phone = temp.RecInService.Phone + Text
					} else {
						//если пользователь нажал инлайн кнопку "Отправить телефон", то записываем его в поле
						temp.RecInService.Phone = temp.RecInService.Phone + update.Message.Contact.PhoneNumber
					}
					users[ID] = temp
					i = 12
					writeUserDataIndex(ID)
					msg.ReplyMarkup = kbrdYNRecInService
					//Выводим записанные данные клиента в виде сообщения
					msg.Text = users[ID].RecInService.Intro +
						"\n" + users[ID].RecInService.Date +
						"\n" + users[ID].RecInService.Time +
						"\n" + users[ID].RecInService.Device +
						"\n" + users[ID].RecInService.Problem +
						"\n" + users[ID].RecInService.Phone +
						"\n" + users[ID].RecInService.Name +
						"\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
				}
			case 12:
				//если вместо ввода данных нажали кнопку главного меню, то не принимаем этот ответ
				if isCommand(Text, btnsMainMenu) {
					msg.Text = "Подтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись, записываясь в сервис"
					msg.ReplyMarkup = kbrdYN
				} else {
					//если НЕ НАЖАТА кнопка "Да, все верно" и отправлено какое-либо сообщение, то отправляем сообщение менеджеру с записанными данными
					msg.ChatID = perm.ManagerID
					msg.Text = "⚡ Запись на сервис:" +
						"\n\n" + users[ID].FirstName +
						"\n" + users[ID].Username +
						"\n" + users[ID].RecInService.Date +
						"\n" + users[ID].RecInService.Time +
						"\n" + users[ID].RecInService.Device +
						"\n" + users[ID].RecInService.Problem +
						"\n" + users[ID].RecInService.Phone +
						"\n" + users[ID].RecInService.Name +
						"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
					if _, err := bot.Send(msg); err != nil {
						panic(err)

					}
					//удаляем из мапы запись с данными пользователя
					delete(users, ID)
					//отправляем сообщение пользователю
					msg.ChatID = ID
					msg.ReplyMarkup = kbrdMain
					msg.Text = "Спасибо, вы записаны на сервис\n "
				}
			case 13: //ветка после нажатия на кнопку Статус заказа
				msg.Text = "Секундочку, проверяю... "
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				tokenRemonline := remonline.TokenRmnln(perm.ApiKey) //TODO не заработал через MustApiKeyRemonline, разобраться
				var phoneForOrderStatus string
				if update.Message.Text != "" {
					//если пользователь отправил телефон текстовым сообщением, то записываем его в переменную phoneForOrderStatus
					phoneForOrderStatus = Text
					msg.Text = remonline.OrderStatus(tokenRemonline, phoneForOrderStatus)
					msg.ReplyMarkup = kbrdMain
				} else if update.Message.Contact != nil {
					//если пользователь нажал инлайн кнопку "Отправить телефон", то записываем его в поле
					phoneForOrderStatus = update.Message.Contact.PhoneNumber
					msg.Text = remonline.OrderStatus(tokenRemonline, phoneForOrderStatus)
					msg.ReplyMarkup = kbrdMain
				} else {
					phoneForOrderStatus = "_"
					msg.Text = remonline.OrderStatus(tokenRemonline, phoneForOrderStatus)
					msg.ReplyMarkup = kbrdMain
				}
				//удаляем из мапы запись с данными пользователя
				delete(users, ID)
			}
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else

		//ветка для обработки команд и нажатия кнопок главного меню при условии, что полученная команда не /start и что user.Index == 0
		if update.Message != nil && update.Message.Text != "/start" && users[update.Message.From.ID].Index == 0 {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// Реакция на нажатия кнопок главного меню
			switch update.Message.Text {
			case "/close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			case perm.OrderStatus:
				msg.Text = perm.EnterPhoneForStatus
				msg.ReplyMarkup = btnContact
				i = 13
				writeUserDataIndex(update.Message.Chat.ID) //используем функцию только для записи в поле users[ID].index значения 13
			case perm.WriteToManager:
				msg.Text = "Нажмите эту кнопку для перехода в чат с менеджером"
				msg.ReplyMarkup = btnURL
			case perm.CallTheOffice:
				msg.Text = perm.CallThisNumber
			case perm.GetCourier:
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
			case perm.RecInService:
				msg.Text = perm.SelectDate
				msg.ReplyMarkup = kbrdDate(t)
			default:
				msg.Text = "Я вас не понимаю. Воспользуйтесь кнопками ниже для общения со мной ⬇️"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		}

		//Обработка нажатия инлайн-кнопок
		if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			ID := update.CallbackQuery.Message.Chat.ID

			//получение значения data после нажатия кнопки с датой
			res := ""
			for _, val := range btnsKbrdDates {
				if update.CallbackQuery.Data == val {
					res = val
				}
			}

			switch update.CallbackQuery.Data {
			case perm.Organization:
				msg.Text = perm.NameOfTheOrganization
				i = 1
				FirstName := userData.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := userData.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUserData(ID, FirstName, Username)
			case perm.NotOrganization:
				msg.Text = perm.Address
				i = 2
				FirstName := userData.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := userData.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUserData(ID, FirstName, Username)
			case perm.Yes:
				// Проверяем, что запись в мапе существует
				if _, ok := users[ID]; ok {
					temp := users[ID]
					//Проверяем, что индекс = 7, т.е. пользователь ответил на все вопросы
					if temp.Index == 7 {
						// Отправляем сообщение менеджеру с данными
						msg.ChatID = perm.ManagerID
						msg.Text = "⚡ Заявка на вызов курьера:" +
							"\n\n" + users[ID].FirstName +
							"\n" + users[ID].Username +
							"\n" + users[ID].GetCourier.Org +
							"\n" + users[ID].GetCourier.Address +
							"\n" + users[ID].GetCourier.Person +
							"\n" + users[ID].GetCourier.Phone +
							"\n" + users[ID].GetCourier.Purpose +
							"\n" + users[ID].GetCourier.Time +
							"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
						//пользователю выводим главное меню и шлем сообщение
						msg.ReplyMarkup = kbrdMain
						msg.ChatID = ID
						msg.Text = perm.Ok
						//удаляем из мапы запись с данными пользователя
						delete(users, ID)
					} else { //При повторном нажатии на кнопку "Да, все верно" и ответах не на все вопросы отработает эта ветка
						msg.Text = "Вы ответили не на все вопросы. Пожалуйста, нажмите кнопку \"Вызвать курьера\" снизу в меню."
						msg.ReplyMarkup = kbrdMain
					}
				} else { //При повторном нажатии на кнопку "Да, все верно" без предварительного ответа на вопросы для вызова курьера отработает эта ветка
					msg.Text = "Ваша заявка на вызов курьера уже была принята. Для оформления новой заявки нажмите кнопку \"Вызвать курьера\" снизу в меню."
				}
			case perm.No:
				//Стираем данные из полей пользователя, перезаписывая их на пустые поля
				i = 0
				FirstName := userData.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := userData.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUserData(ID, FirstName, Username)
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
			case res:
				i = 0
				FirstName := userData.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := userData.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUserData(ID, FirstName, Username)
				//записываем дату в поле Date
				temp := users[ID]
				temp.RecInService.Date = temp.RecInService.Date + update.CallbackQuery.Data
				users[ID] = temp
				msg.Text = "Выберете удобное время:"
				msg.ReplyMarkup = kbrdBeforeAfterLunch //рисуем кнопки "до/после обеда"
			case "beforeLunch":
				i = 8
				//записываем время в поле Time
				temp := users[ID]
				temp.RecInService.Time = temp.RecInService.Time + "До обеда"
				users[ID] = temp
				writeUserDataIndex(ID)
				msg.Text = "Какая у вас проблема? (например, разбился экран телефона)"
			case "afterLunch":
				i = 8
				//записываем время в поле Time
				temp := users[ID]
				temp.RecInService.Time = temp.RecInService.Time + "После обеда"
				users[ID] = temp
				writeUserDataIndex(ID)
				msg.Text = "Какая у вас проблема? Например, разбился экран телефона."
			case "yes":
				// Проверяем, что запись в мапе существует
				if _, ok := users[ID]; ok {
					temp := users[ID]
					//Проверяем, что индекс = 7, т.е. пользователь ответил на все вопросы
					if temp.Index == 12 {
						// Отправляем сообщение менеджеру с данными
						msg.ChatID = perm.ManagerID
						msg.Text = "⚡ Запись в Сервис:" +
							"\n\n" + users[ID].FirstName +
							"\n" + users[ID].Username +
							"\n" + users[ID].RecInService.Date +
							"\n" + users[ID].RecInService.Time +
							"\n" + users[ID].RecInService.Device +
							"\n" + users[ID].RecInService.Problem +
							"\n" + users[ID].RecInService.Phone +
							"\n" + users[ID].RecInService.Name +
							"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
						//пользователю выводим главное меню и шлем сообщение
						msg.ReplyMarkup = kbrdMain
						msg.ChatID = ID
						msg.Text = "Спасибо, вы записаны."
						//удаляем из мапы запись с данными пользователя
						delete(users, ID)
					} else { //При повторном нажатии на кнопку "Да, все верно" и ответах не на все вопросы отработает эта ветка
						msg.Text = "Вы ответили не на все вопросы. Пожалуйста, нажмите кнопку \"Запись в сервис\" снизу в меню."
					}
				} else { //При повторном нажатии на кнопку "Да, все верно" без предварительного ответа на вопросы для вызова курьера отработает эта ветка
					msg.Text = "Вы уже записаны. Для новой записи нажмите кнопку \"Запись в сервис\" снизу в меню."
					msg.ReplyMarkup = kbrdMain
				}
			case "no":
				//Стираем данные из полей пользователя, перезаписывая их на пустые поля
				i = 0
				FirstName := userData.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := userData.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUserData(ID, FirstName, Username)
				msg.Text = perm.SelectDate
				msg.ReplyMarkup = kbrdDate(t)

			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)

			}
		} else
		//ветка для обработки команды "/start"
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.Text = perm.StartMsg
			msg.ReplyMarkup = kbrdMain
			//удаляем данные клиента
			delete(users, update.Message.Chat.ID)

			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		}

		log.Printf("Содержимое мапы с данными пользователей users map[int64]user:::::%+v\n", users)
	}
}
