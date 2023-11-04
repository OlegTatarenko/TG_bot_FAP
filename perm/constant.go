package perm

// Контактные данные офиса
const (
	Token       = "6459859641:AAEyvBu87_bFrIjsJumyv4KA0gpF1GR6d94"
	OfficePhone = "+79045340560"              //Моб телефон менеджера в офисе
	ChatURL     = "https://t.me/ultrafilaret" //Ссылка на username в телеге, привязанный к моб телефону менеджера.
)

// Названия кнопок
const (
	GetCourier      = "🏃 Вызвать курьера"
	OrderStatus     = "❓ Статус заказа"
	WriteToManager  = "✏ Написать менеджеру"
	CallTheOffice   = "📞 Звонок в офис"
	RecordInService = "🔧 Запись в Сервис"
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

// Срез для записи информации клиента при вызове курьера
var Form = []string{
	"Спасибо! Вы ввели следующие данные: ",
	"Учреждение (если применимо): ",
	"Адрес, где забрать: ",
	"Контактное лицо: ",
	"Контактный телефон: ",
	"Цель вызова курьера: ",
	"Дата, время приезда курьера: ",
}

// Хеш-таблица для хранения данных ползвателей
var Forms map[int64]UserData

// Структура для записи данных пользоватея при вызове курера или при записи на сервис
type UserData struct {
	UserID           int64
	DataGetCourier   Blank
	DataRecInService string
	Index            int
}

type Blank struct {
	Intro   string
	Org     string
	Address string
	Person  string
	Phone   string
	Purpose string
	Time    string
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
