package main

import (
	"TG_bot_FAP/perm"
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
	kbrdDates = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Yes, "day1"),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, "day2"),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, "day3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(perm.Yes, "day4"),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, "day5"),
			tgbotapi.NewInlineKeyboardButtonData(perm.No, "day6"),
		))
)

// Счетчик для вывода инлайн-кнопок и записи данных пользователя
var i = 0

// forms мапа для хранения данных пользователей
var forms = make(map[int64]UserData)

// UserData структура для записи данных пользователя при вызове курьера или при записи на сервис
type UserData struct {
	FirstName        string
	Username         string
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

// TODO: Добавить поле для ссылки на чат
var ud = UserData{
	"Имя: ",
	"Username: ",
	b,
	"Вызов курьера",
	0,
}

var b = Blank{
	"Спасибо! Вы ввели следующие данные: ",
	"Учреждение (если применимо): ",
	"Адрес, где забрать: ",
	"Контактное лицо: ",
	"Контактный телефон: ",
	"Цель вызова курьера: ",
	"Дата, время приезда курьера: ",
}

var t = time.Now()

// writeUDStart функция для создания записи пользователя в мапе с ID пользователя
func writeUDStart(ID int64, FirstName, Username string) {
	if _, ok := forms[ID]; !ok {
		forms[ID] = ud    //Если нет ключа, равного ID, то создаем элемент мапы с ключом = ID
		temp := forms[ID] //переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
		//т.к. изменение значения поля через forms[ID].index = i языком не предусмотрено
		temp.FirstName = FirstName //Записываем в поле имя пользователя
		temp.Username = Username   //Записываем в поле Username пользователя
		temp.Index = i             // Меняем индекс у пользователя
		forms[ID] = temp           //Изменяем запись в мапе
	} else {
		// меняем индекс у пользователя
		temp := forms[ID]
		temp.Index = i
		forms[ID] = temp
	}
}

// writeUDIndex функция для изменения индекса в данных пользователя
func writeUDIndex(ID int64) {
	//переменная с копией структуры, чтобы изменить значение поля структуры внутри мапы,
	//т.к. изменение значения поля через forms[ID].index = i языком не предусмотрено
	//счетчик записываем в index
	temp := forms[ID]
	temp.Index = i
	forms[ID] = temp
}

//TODO: Тест - запись в субботу в между 14:00 и 14:30

// getKeyboardDate формирует клавиатуру из 6 инлайн кнопок, для чего создает слайс длиной 6 из строк,
// где каждая строка - дата для инлайн кнопки на запись в сервис
func getKeyboardDate(t time.Time) tgbotapi.InlineKeyboardMarkup {
	sliceDate := make([]string, 6, 6)
	j := 0                   // индексы сдвига даты
	hh, mm, _ := t.Clock()   //берем часы и минуты из времени нажатия на кнопку
	cntrlTimeSatHH := 14     // контрольное время для субботы 14:00 - для сравнения значения часов
	cntrlTimeSatMM := 0      // контрольное время для субботы 14:00 - для сравнения значения минут
	cntrlTimeWorkdayHH := 16 //контрольное время для будних дней 16:30 - для сравнения значения часов
	cntrlTimeWorkdayMM := 30 //контрольное время для будних дней 16:30 - для сравнения значения минут
	for k := 0; k < len(sliceDate); k++ {
		l := 0 // индексы сдвига даты
		//если время нажатия на кнопку > 16:30 (cntrlTimeWorkday) или день нажатия на кнопку - суббота и время > 14:00 (cntrlTimeSat), то
		//увеличиваем l на единицу, чтобы пропустить сегодняшнюю дату
		if hh > cntrlTimeWorkdayHH || (hh == cntrlTimeWorkdayHH && mm > cntrlTimeWorkdayMM) || (t.Weekday() == time.Saturday && (hh > cntrlTimeSatHH || hh == cntrlTimeSatHH && mm > cntrlTimeSatMM)) {
			l = 1
		}

		datesRecInService := t.Add(time.Hour * 24 * time.Duration(k+j+l)) //добавляем к текущей дате 24*(k+j+l) часов
		weekDay := datesRecInService.Weekday()                            //определяем для даты день недели

		//если день недели - воскресенье, т.е. выходной, то j=1 , чтобы пропустить дату воскресенья
		if weekDay == time.Sunday {
			j = 1
			datesRecInService = t.Add(time.Hour * 24 * time.Duration(k+j+l))
		}
		sliceDate[k] = datesRecInService.Format("02.01.2006") //форматируем дату в строку и записываем в слайс
	}
	kbrdDates = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[0], "day1"),
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[1], "day2"),
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[2], "day3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[3], "day4"),
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[4], "day5"),
			tgbotapi.NewInlineKeyboardButtonData(sliceDate[5], "day6"),
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

		//TODO: Добавить условие, чтобы не принимал нажатие кнопок главного меню в качестве ответов в этой ветке

		// Опрос клиента после нажатия кнопки Вызвать курьера, т.е. при условии, что UserData.Index > 0
		if update.Message != nil && forms[update.Message.From.ID].Index > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			ID := update.Message.Chat.ID
			Text := update.Message.Text
			switch forms[update.Message.From.ID].Index {
			case 1:
				//записываем полученное наименование организации в поле Org
				temp := forms[ID]
				temp.DataGetCourier.Org = temp.DataGetCourier.Org + Text
				forms[ID] = temp
				i = 2
				writeUDIndex(ID)
				msg.Text = perm.Adress
			case 2:
				//записываем полученный адрес в поле Address
				temp := forms[ID]
				temp.DataGetCourier.Address = temp.DataGetCourier.Address + Text
				forms[ID] = temp
				i = 3
				writeUDIndex(ID)
				msg.Text = "Укажите имя контактного лица"
			case 3:
				//записываем полученное имя в поле Person
				temp := forms[ID]
				temp.DataGetCourier.Person = temp.DataGetCourier.Person + Text
				forms[ID] = temp
				i = 4
				writeUDIndex(ID)
				msg.ReplyMarkup = btnContact
				msg.Text = "Введите номер телефона или нажмите кнопку, чтобы отправить номер телефона, к которому привязан ваш аккаунт телеграма"

			case 4:
				//записываем полученный телефон в поле Phone
				temp := forms[ID]
				if update.Message.Text != "" {
					//если пользователь отправил телефон текстовым сообщением, то записываем его в поле
					temp.DataGetCourier.Phone = temp.DataGetCourier.Phone + Text
				} else {
					//если пользователь нажал инлайн кнопку "Отправить телефон", то записываем его в поле
					temp.DataGetCourier.Phone = temp.DataGetCourier.Phone + update.Message.Contact.PhoneNumber
				}

				forms[ID] = temp
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
				temp := forms[ID]
				temp.DataGetCourier.Purpose = temp.DataGetCourier.Purpose + Text
				forms[ID] = temp
				i = 6
				writeUDIndex(ID)
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Укажите удобную дату и время приезда курьера"
			case 6:
				//записываем полученную дату и время в поле Time
				temp := forms[ID]
				temp.DataGetCourier.Time = temp.DataGetCourier.Time + Text
				forms[ID] = temp
				i = 7
				writeUDIndex(ID)
				msg.ReplyMarkup = kbrdYN
				//Выводим записанные данные клиента в виде сообщения
				msg.Text = forms[ID].DataGetCourier.Intro +
					"\n" + forms[ID].DataGetCourier.Org +
					"\n" + forms[ID].DataGetCourier.Address +
					"\n" + forms[ID].DataGetCourier.Person +
					"\n" + forms[ID].DataGetCourier.Phone +
					"\n" + forms[ID].DataGetCourier.Purpose +
					"\n" + forms[ID].DataGetCourier.Time +
					"\n\nПодтвердите, если все верно, или нажмите \"Исправить\", если вы ошиблись."
			case 7:
				//если НЕ НАЖАТА кнопка "Да, все верно" и отправлено какое-либо сообщение, то отправляем сообщение менеджеру с записанными данными
				msg.ChatID = perm.ManagerID
				msg.Text = "⚡ Заявка на вызов курьера:" +
					"\n\n" + forms[ID].FirstName +
					"\n" + forms[ID].Username +
					"\n" + forms[ID].DataGetCourier.Org +
					"\n" + forms[ID].DataGetCourier.Address +
					"\n" + forms[ID].DataGetCourier.Person +
					"\n" + forms[ID].DataGetCourier.Phone +
					"\n" + forms[ID].DataGetCourier.Purpose +
					"\n" + forms[ID].DataGetCourier.Time +
					"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
				if _, err := bot.Send(msg); err != nil {
					panic(err)

				}
				//удаляем из мапы запись с данными пользователя
				delete(forms, ID)
				//отправляем сообщение пользователю
				msg.ChatID = ID
				msg.ReplyMarkup = kbrdMain
				msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки.\n "
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

			switch update.CallbackQuery.Data {
			case perm.Organization:
				msg.Text = perm.NameOfTheOrganization
				i = 1
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
			case perm.NotOrganization:
				msg.Text = perm.Adress
				i = 2
				FirstName := ud.FirstName + update.CallbackQuery.From.FirstName    //Записываем в поле имя пользователя
				Username := ud.Username + "@" + update.CallbackQuery.From.UserName //Записываем в поле Username пользователя
				writeUDStart(ID, FirstName, Username)
			case perm.Yes:
				// Проверяем, что запись в мапе существует
				if _, ok := forms[ID]; ok {
					temp := forms[ID]
					//Проверяем, что индекс = 7, т.е. пользователь ответил на все вопросы
					if temp.Index == 7 {
						// Отправляем сообщение менеджеру с данными
						msg.ChatID = perm.ManagerID
						msg.Text = "⚡ Заявка на вызов курьера:" +
							"\n\n" + forms[ID].FirstName +
							"\n" + forms[ID].Username +
							"\n" + forms[ID].DataGetCourier.Org +
							"\n" + forms[ID].DataGetCourier.Address +
							"\n" + forms[ID].DataGetCourier.Person +
							"\n" + forms[ID].DataGetCourier.Phone +
							"\n" + forms[ID].DataGetCourier.Purpose +
							"\n" + forms[ID].DataGetCourier.Time +
							"\n\n" + "Свяжитесь с клиентом для подтверждения заявки"
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
						//пользователю выводим главное меню и шлем сообщение
						msg.ReplyMarkup = kbrdMain
						msg.ChatID = ID
						msg.Text = "Спасибо, ваша заявка принята. В ближайшее время с вами свяжется менеджер по указанному телефону для подтверждения заявки."
						//удаляем из мапы запись с данными пользователя
						delete(forms, ID)
					} else { //При повторном нажатии на кнопку "Да, все верно" и ответах не на все вопросы отработает эта ветка
						msg.Text = "Вы ответили не на все вопросы. Пожалуйста, нажмите кнопку \"Вызвать курьера\" снизу в меню."
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
			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)

			}
		}

		fmt.Println("\v", forms)
		fmt.Println("\v", ud)
	}
}
