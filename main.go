package main

import (
	"TG_bot_FAP/perm"
	"TG_bot_FAP/remonline"
	"fmt"
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

//TODO требует рестарта после перехода в чат менеджеру и возврата к боту - исправить

// Счетчик для вывода инлайн-кнопок и записи данных пользователя
var i = 0

// users мапа для хранения данных пользователей
var users = make(map[int64]user)

// sliceBtnString срез для хранения значений "text" кнопок в клавиатуре getKeyboardDate
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

// blankGetCourier структура для записи данных пользователя при вызове курьера
type blankRecInService struct {
	Intro   string
	Date    string
	Time    string
	Device  string
	Problem string
	Phone   string
	Name    string
}

var dataGetCourier = blankGetCourier{
	"Спасибо! Вы ввели следующие данные: ",
	"Учреждение (если применимо): ",
	"Адрес, где забрать: ",
	"Имя: ",
	"Номер телефона: ",
	"Цель вызова курьера: ",
	"Дата, время приезда курьера: ",
}
var dataRecInService = blankRecInService{
	"Спасибо! Вы ввели следующие данные: ",
	"Дата: ",
	"Время: ",
	"Устройство: ",
	"Проблема: ",
	"Номер телефона: ",
	"Имя: ",
}

// TODO: Добавить поле для ссылки на чат
var ud = user{
	"Name: ",
	"Username: ",
	dataGetCourier,
	dataRecInService,
	0,
}

var t = time.Now()

// writeUDStart функция для создания записи пользователя в мапе с ID пользователя
func writeUDStart(ID int64, FirstName, Username string) {
	if _, ok := users[ID]; !ok {
		users[ID] = ud    //Если нет ключа, равного ID, то создаем элемент мапы с ключом = ID
		temp := users[ID] //переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
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

// writeUDIndex функция для изменения индекса в данных пользователя
func writeUDIndex(ID int64) {
	//переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
	//т.к. изменение значения поля через users[ID].index = i языком не предусмотрено
	//счетчик записываем в index
	temp := users[ID]
	temp.Index = i
	users[ID] = temp
}

//TODO: Тест - запись в субботу в между 14:00 и 14:30
//TODO: Контрольное время в константу

// getKeyboardDate формирует клавиатуру из 6 инлайн кнопок, для чего создает слайс длиной 6 из строк,
// где каждая строка - дата для инлайн кнопки на запись в сервис
func getKeyboardDate(t time.Time) tgbotapi.InlineKeyboardMarkup {
	j := 0                         // индексы сдвига даты
	hh, mm, _ := t.Clock()         //берем часы и минуты из времени нажатия на кнопку, т.е. из текущего времени
	timeNOW := hh*60 + mm          //текущее время в минутах с начала суток
	cntrlTimeSat := 14 * 60        // контрольное время для субботы 14:00 - в минутах с начала суток
	cntrlTimeWorkday := 16*60 + 30 //контрольное время для будних дней 16:30 - в минутах с начала суток
	for k := 0; k < len(btnsKbrdDates); k++ {
		l := 0 // индексы сдвига даты
		//если время нажатия на кнопку > 16:30 (cntrlTimeWorkday) или день нажатия на кнопку - суббота и время > 14:00 (cntrlTimeSat), то
		//увеличиваем l на единицу, чтобы пропустить сегодняшнюю дату
		if timeNOW > cntrlTimeWorkday || (t.Weekday() == time.Saturday && timeNOW > cntrlTimeSat) {
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

		//TODO: Добавить условие, чтобы не принимал нажатие кнопок главного меню в качестве ответов в этой ветке или убирать кнопки главного меню

		// Опрос клиента после нажатия кнопок Вызвать курьера, Запись в сервис или Статус заказа, т.е. при условии, что user.Index > 0
		if update.Message != nil && users[update.Message.From.ID].Index > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			ID := update.Message.Chat.ID
			Text := update.Message.Text
			switch users[update.Message.From.ID].Index {
			case 1:
				//записываем полученное наименование организации в поле Org
				temp := users[ID]
				temp.GetCourier.Org = temp.GetCourier.Org + Text
				users[ID] = temp
				i = 2
				writeUDIndex(ID)
				msg.Text = perm.Address
			case 2:
				//записываем полученный адрес в поле Address
				temp := users[ID]
				temp.GetCourier.Address = temp.GetCourier.Address + Text
				users[ID] = temp
				i = 3
				writeUDIndex(ID)
				msg.Text = "Укажите имя контактного лица"
			case 3:
				//записываем полученное имя в поле Person
				temp := users[ID]
				temp.GetCourier.Person = temp.GetCourier.Person + Text
				users[ID] = temp
				i = 4
				writeUDIndex(ID)
				msg.ReplyMarkup = btnContact
				msg.Text = "Введите номер телефона или нажмите кнопку, чтобы отправить свой номер телефона, к которому привязан ваш аккаунт телеграма"
			case 4:
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
				writeUDIndex(ID)
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Для чего вызываете курьера? Например, часто указывают:" +
					"\n- забрать картриджи, заправить, вернуть." +
					"\n- купить 1 новый картридж ce285a и выставить счет." +
					"\n- после заправки привезти акт сверки и новый счет." +
					"\n- забрать подписанный договор." +
					"\n...или ваш вариант"
			case 5:
				//записываем полученную цель вызова в поле Purpose
				temp := users[ID]
				temp.GetCourier.Purpose = temp.GetCourier.Purpose + Text
				users[ID] = temp
				i = 6
				writeUDIndex(ID)
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Укажите удобную дату и время приезда курьера"
			case 6:
				//записываем полученную дату и время в поле Time
				temp := users[ID]
				temp.GetCourier.Time = temp.GetCourier.Time + Text
				users[ID] = temp
				i = 7
				writeUDIndex(ID)
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
			case 7:
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
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки.\n "
			case 8:
				// TODO записать проблему в структуру
				//записываем проблему в поле Problem
				temp := users[ID]
				temp.RecInService.Problem = temp.RecInService.Problem + Text
				users[ID] = temp
				i = 9
				writeUDIndex(ID)
				msg.Text = "Какое у вас устройство? Например, ноутбук ACER модель ABC."
			case 9:
				// TODO записать устройство в структуру
				//записываем устройство в поле Device
				temp := users[ID]
				temp.RecInService.Device = temp.RecInService.Device + Text
				users[ID] = temp
				i = 10
				writeUDIndex(ID)
				msg.Text = "Укажите ваше имя"
			case 10:
				// TODO записать имя в структуру
				//записываем имя в поле Name
				temp := users[ID]
				temp.RecInService.Name = temp.RecInService.Name + Text
				users[ID] = temp
				i = 11
				writeUDIndex(ID)
				msg.ReplyMarkup = btnContact
				msg.Text = "Введите номер телефона или нажмите кнопку, чтобы отправить номер телефона, к которому привязан ваш аккаунт телеграма"
			case 11:
				//TODO записываем полученный телефон в поле Phone
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
				writeUDIndex(ID)
				msg.ReplyMarkup = kbrdYNRecInService
				//TODO Выводим записанные данные клиента в виде сообщения
				msg.Text = users[ID].RecInService.Intro +
					"\n" + users[ID].RecInService.Date +
					"\n" + users[ID].RecInService.Time +
					"\n" + users[ID].RecInService.Device +
					"\n" + users[ID].RecInService.Problem +
					"\n" + users[ID].RecInService.Phone +
					"\n" + users[ID].RecInService.Name +
					"\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
			case 12:
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
			case 13: //ветка после нажатия на кнопку Статус заказа
				var phoneForOrderStatus string
				if update.Message.Text != "" {
					//если пользователь отправил телефон текстовым сообщением, то записываем его в переменную phoneForOrderStatus
					phoneForOrderStatus = Text
					tokenRemonline := remonline.Token(perm.ApiKey)
					msg.Text = remonline.OrderStatus(tokenRemonline, phoneForOrderStatus)
					msg.ReplyMarkup = kbrdMain
				} else {
					//если пользователь нажал инлайн кнопку "Отправить телефон", то записываем его в поле
					phoneForOrderStatus = update.Message.Contact.PhoneNumber
					tokenRemonline := remonline.Token(perm.ApiKey)
					msg.Text = remonline.OrderStatus(tokenRemonline, phoneForOrderStatus)
					msg.ReplyMarkup = kbrdMain
				}
				//удаляем из мапы запись с данными пользователя
				delete(users, ID)
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
				msg.Text = "Введите номер телефона, указанный в заказе, не менее 10 цифр в формате 9123456789. Для отправки своего номера можете нажать кнопку ниже."
				msg.ReplyMarkup = btnContact
				i = 13
				writeUDIndex(update.Message.Chat.ID) //используем функцию только для записи в поле users[ID].index значения 13
			case perm.WriteToManager:
				msg.Text = "Нажмите эту кнопку для перехода в чат с менеджером"
				msg.ReplyMarkup = btnURL
			case perm.CallTheOffice:
				msg.Text = perm.CallThisNumber
			case perm.GetCourier:
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
			case perm.RecInService:
				msg.Text = "Выберите удобную дату"
				msg.ReplyMarkup = getKeyboardDate(t)
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
			ID := update.CallbackQuery.Message.Chat.ID

			//чудо-костыль для получения значения data после нажатия кнопки с датой
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
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
			case perm.NotOrganization:
				msg.Text = perm.Address
				i = 2
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
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
						msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки."
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
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
				msg.Text = perm.AreYouOrg
				msg.ReplyMarkup = kbrdYNOrg
			case res:
				// TODO записать дату в структуру
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
				//записываем дату в поле Date
				temp := users[ID]
				temp.RecInService.Date = temp.RecInService.Date + update.CallbackQuery.Data
				users[ID] = temp
				msg.Text = "Выберете удобное время:"   // TODO прописать сообщение в константы
				msg.ReplyMarkup = kbrdBeforeAfterLunch //рисуем кнопки "до/после обеда"
			case "beforeLunch":
				// TODO записать время в структуру
				i = 8
				//записываем время в поле Time
				temp := users[ID]
				temp.RecInService.Time = temp.RecInService.Time + "До обеда"
				users[ID] = temp
				writeUDIndex(ID)
				msg.Text = "Какая у вас проблема? (например, разбился экран телефона)" // TODO прописать сообщение в константы
			case "afterLunch":
				// TODO записать время в структуру
				i = 8
				//записываем время в поле Time
				temp := users[ID]
				temp.RecInService.Time = temp.RecInService.Time + "После обеда"
				users[ID] = temp
				writeUDIndex(ID)
				msg.Text = "Какая у вас проблема? Например, разбился экран телефона." // TODO прописать сообщение в константы
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
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
				msg.Text = "Выберите удобную дату. Если дата попадает на праздничный день, лучше написать менеджеру или позвонить в офис для уточнения графика работы."
				msg.ReplyMarkup = getKeyboardDate(t)

			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)

			}
		}
		fmt.Println("\v", users)
		fmt.Println("\v", ud)
	}
}
