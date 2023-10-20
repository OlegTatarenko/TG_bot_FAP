package perm

// Контактные данные офиса
const (
	OfficePhone = "+79045340560"              //Моб телефон менеджера в офисе
	ChatURL     = "https://t.me/ultrafilaret" //Ссылка на никнейм в телеге, привязанный к моб телефону менеджера.
)

// Названия кнопок
const (
	GetCourier      = "🏃 Вызвать курьера"
	OrderStatus     = "❓ Статус заказа"
	WriteToManager  = "✏ Написать менеджеру"
	CallTheOffice   = "📞 Позвонить в офис"
	RecordInService = "🔧 Записаться в Сервис"
	Organization    = "Я представляю учреждение 🏢"
	NotOrganization = "Нет, я не учреждение 🙂"
	Back            = "Назад 🔙"
	GoToChat        = "➡ Нажмите для перехода в чат ⬅"
	CallThisNumber  = "☎ офиса : " + OfficePhone
)

// Ответы бота в чат
const (
	InDev                 = "Извините, эта функция еще в разработке🤷‍♂️"
	NameOfTheOrganization = "Укажите наименование учреждения"
)
