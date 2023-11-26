package perm

// Контактные данные офиса
const (
	Token             = "6459859641:AAEyvBu87_bFrIjsJumyv4KA0gpF1GR6d94"
	OfficePhone       = "+79045340560"              //Моб телефон менеджера в офисе
	ChatURL           = "https://t.me/ultrafilaret" //Ссылка на username в телеге, привязанный к моб телефону менеджера
	ManagerID   int64 = 47526499                    // MenedgerID переменная с id офис-менеджера в телеграм
)

// Названия кнопок
const (
	GetCourier      = "🏃 Вызвать курьера"
	OrderStatus     = "❓ Статус заказа"
	WriteToManager  = "✏ Написать менеджеру"
	CallTheOffice   = "📞 Звонок в офис"
	RecInService    = "🔧 Запись в Сервис"
	Organization    = "Да 🏢"
	NotOrganization = "Нет 🙂"
	Back            = "Назад 🔙"
	GoToChat        = "▶ Нажмите сюда ◀"
	CallThisNumber  = "Телефон офиса : " + OfficePhone
	Adress          = "Укажите полный адрес, где забрать"
	Contact         = "▶ Отправить номер ◀"
	Yes             = "Да, все верно 👍"
	No              = "Исправить❗"
	AreYouOrg       = "Вы представляете учреждение?"
)

// Ответы бота в чат
const (
	InDev                 = "Извините, эта функция еще в разработке🤷‍♂️"
	NameOfTheOrganization = "Укажите наименование учреждения"
)
