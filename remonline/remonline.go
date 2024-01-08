package remonline

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"tgbot_Smartset/perm"
	"time"
)

// MustApiKeyRemonline - функция для получения ApiKey remonline через флаг -remonline-apikey
// для запуска из командной строки необходимо:
//   - go build (собираем exe-файл, если его еще нет, если есть - пропускаем эту команду)
//   - ./TG_bot_FAP -tgbot-token 'значение токена тг бота' -remonline-apiKey 'значение apikey' (запускаем exe-файл с флагами, указывая значение токена и apikey)
func MustApiKeyRemonline() string {
	apiKey := flag.String(
		"remonline-apiKey",
		"",
		"apiKey для доступа к remonline / apiKey for access to remonline",
	)

	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey не указан / apiKey is not specified")
	}

	return *apiKey
}

// TokenRmnln - функция для получения токена от Ремонлайн по apiKey Смартсет, токен действителен сутки с момента получения
func TokenRmnln(apiKey string) string {
	//получаем токен, используя в запросе к API Ремонлайна apiKey Смартсет (https://remonline.app/docs/api/#apisection210_546)
	url := "https://api.remonline.app/token/new?api_key=" + apiKey

	//у клиента по умолчанию не установлен тайм-аут, пропишем его;
	//если удаленный сервер не отвечает, то клиент без установленного тайм-аута будет ожидать ответ бесконечно?
	var myClient = &http.Client{Timeout: 10 * time.Second}
	//создаем новый запрос
	req, err := http.NewRequest("GET", url, nil)
	//получаем данные
	resp, err := myClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	//закрываем тело ответа, чтобы избежать утечки ресурсов
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)
	//читаем содержимое тела (из поля Body структуры Response) с помощью ReadAll в срез байт
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//проверка json на правильность
	if !json.Valid(body) {
		log.Fatalln("invalid json!") // вывод: invalid json!
	}
	//структура и переменная для сохранения токена
	type Tkn struct {
		Token string `json:"token"`
	}
	var T Tkn
	//декодируем полученный json (содержимое переменной body - срез байт), записываем токен в переменную T
	if err = json.Unmarshal(body, &T); err != nil {
		log.Fatalln(err)
		//return
	}
	// этот вариант лучше на случай некорректных данных в json, но он выдает EOF, не разобрался
	//err = json.NewDecoder(resp.Body).Decode(&T)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//возвращаем значение токен
	return T.Token
}

// IsValidPhone - функция, проверяющая корректность ввода номера телефона: допустимы цифры, знак +, не менее 11 символов
func IsValidPhone(phone string) bool {
	res := false
	//приводим строку к срезу рун, если его длина < 11, то номер введен некорректно
	if len([]rune(phone)) < 11 {
		return res
	} else {
		for _, val := range phone {
			//допустимы только цифры 0-9 (коды ASCII в десятичной системе сч-я: 48-57) и знак + (код ASCII в десятичной системе сч-я: 43)
			if (47 < val && val < 58) || val == 43 {
				res = true
			} else {
				res = false
				break
			}
		}
		return res
	}
}

// titleOfStatus - функция для выбора варианта заголовков в ответе бота перед статусом заказа/ов
func titleOfStatus(nmbrOfOrders int, nmbrOfOutputOrders int) string {
	var title string
	switch nmbrOfOrders {
	case 1:
		title = "Статус заказа ℹ\n"
	case 2:
		title = "Найдено 2 заказа ℹ\n"
	case 3:
		title = "Найдено 3 заказа ℹ\n"
	case 4:
		title = "Найдено 4 заказа ℹ\n"
	case 5:
		title = "Найдено 5 заказов ℹ\n"
	case 6:
		title = "Найдено 6 заказов ℹ\n"
	default:
		title = "Найдено заказов: " + strconv.Itoa(nmbrOfOrders) + "\nПокажу последние " + strconv.Itoa(nmbrOfOutputOrders) + " ℹ\n"
	}
	return title
}

// descriptionOfStatus - функция для присвоения каждой группе статуса словесного описания
func descriptionOfStatus(statusGroup int) string {
	var status string
	switch statusGroup {
	case 0:
		status = "Пользовательский"
	case 1:
		status = "Новый"
	case 2:
		status = "На исполнении"
	case 3:
		status = "Отложен"
	case 4:
		status = "Выполнен"
	case 5:
		status = "Доставка"
	case 6:
		status = "Закрыт успешно"
	case 7:
		status = "Закрыт неуспешно"
	default:
		status = "Нет информации"
	}
	return status
}

// OrderStatus - функция для получения данных статуса заказа по номеру телефона клиента
func OrderStatus(token, phone string) string {
	output := ""
	if !IsValidPhone(phone) {
		output = perm.NotCorrectPhone
	} else {
		//получаем из списка клиентов данные по конкретному клиенту,
		// используя в запросе к API Ремонлайна номер телефона клиента и сортировку по убыванию номеров заказов (https://remonline.app/docs/api/#apisection212)
		url := "https://api.remonline.app/order/?token=" + token + "&client_phones[]=" + phone + "&sort_dir=desc" // + "&statuses[]=[1],[6]" не разобрался, как прописать фильтр по статусу заказа

		//у клиента по умолчанию не установлен тайм-аут, пропишем его;
		//если удаленный сервер не отвечает, то клиент без установленного тайм-аута будет ожидать ответ бесконечно?
		var myClient = &http.Client{Timeout: 10 * time.Second}
		//создаем новый запрос
		req, err := http.NewRequest("GET", url, nil)
		//получаем список заказов, отфильтрованный по номеру телефона клиента и убыванию номеров заказов
		resp, err := myClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		//закрываем тело ответа, чтобы избежать утечки ресурсов
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(resp.Body)
		//читаем содержимое тела (из поля Body структуры Response) с помощью ReadAll в срез байт
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//проверка json на правильность
		if !json.Valid(body) {
			log.Fatalln("invalid json!") // вывод: invalid json!
		}
		//структура и переменная для сохранения данных, сформировал через https://mholt.github.io/json-to-go/
		type Order struct {
			Count int `json:"-"`
			Data  []struct {
				ID     int     `json:"id"`
				Brand  string  `json:"brand"`
				Model  string  `json:"model"`
				Price  float64 `json:"price"`
				Payed  float64 `json:"-"`
				Resume string  `json:"-"`
				Urgent bool    `json:"-"`
				Serial string  `json:"-"`
				Client struct {
					ID                int      `json:"-"`
					Phone             []string `json:"-"`
					Address           string   `json:"-"`
					Name              string   `json:"-"`
					Email             string   `json:"-"`
					ModifiedAt        int64    `json:"-"`
					Notes             string   `json:"-"`
					Supplier          bool     `json:"-"`
					Juridical         bool     `json:"-"`
					Conflicted        bool     `json:"-"`
					DiscountCode      string   `json:"-"`
					DiscountGoods     int      `json:"-"`
					DiscountServices  int      `json:"-"`
					DiscountMaterials int      `json:"-"`
					CustomFields      struct {
						Num1 string `json:"-"`
					} `json:"-"`
					MarketingSource struct {
						ID    int    `json:"-"`
						Title string `json:"-"`
					} `json:"-"`
				} `json:"-"`
				Status struct {
					ID    int    `json:"-"`
					Name  string `json:"-"`
					Group int    `json:"group"`
					Color string `json:"-"`
				} `json:"status"`
				DoneAt      int64  `json:"-"`
				Overdue     bool   `json:"-"`
				EngineerID  int    `json:"-"`
				ManagerID   int    `json:"-"`
				BranchID    int    `json:"-"`
				Appearance  string `json:"-"`
				CreatedByID int    `json:"-"`
				OrderType   struct {
					ID    int    `json:"-"`
					Title string `json:"-"`
				} `json:"-"`
				Parts []struct {
					ID             int     `json:"-"`
					EngineerID     int     `json:"-"`
					Title          string  `json:"-"`
					Cost           float64 `json:"-"`
					Price          float64 `json:"-"`
					DiscountValue  int     `json:"-"`
					Amount         float64 `json:"-"`
					Warranty       int     `json:"-"`
					WarrantyPeriod int     `json:"-"`
				} `json:"-"`
				Operations []struct {
					ID             int    `json:"-"`
					EngineerID     int    `json:"-"`
					Title          string `json:"-"`
					Cost           int    `json:"-"`
					Price          int    `json:"-"`
					DiscountValue  int    `json:"-"`
					Amount         int    `json:"-"`
					Warranty       int    `json:"-"`
					WarrantyPeriod int    `json:"-"`
				} `json:"-"`
				Attachments []struct {
					CreatedByID int    `json:"-"`
					CreatedAt   int64  `json:"-"`
					URL         string `json:"-"`
					Filename    string `json:"-"`
				} `json:"-"`
				CreatedAt    int64  `json:"-"`
				ScheduledFor any    `json:"-"`
				ClosedAt     int64  `json:"-"`
				ModifiedAt   int64  `json:"-"`
				Packagelist  string `json:"-"`
				KindofGood   string `json:"-"`
				Malfunction  string `json:"-"`
				IDLabel      string `json:"-"`
				ClosedByID   int    `json:"-"`
				CustomFields struct {
					Num1 string `json:"1"`
				} `json:"-"`
				WarrantyDate    int64  `json:"-"`
				ManagerNotes    string `json:"-"`
				EstimatedCost   int    `json:"-"`
				EngineerNotes   string `json:"-"`
				WarrantyGranted bool   `json:"-"`
				EstimatedDoneAt int64  `json:"-"`
			} `json:"data"`
			Page    int  `json:"-"`
			Success bool `json:"-"`
		}
		var User Order
		//декодируем полученный json (содержимое переменной body - срез байт), записываем данные в переменную User
		if err = json.Unmarshal(body, &User); err != nil {
			log.Fatalln(err)
			//return
		}
		//если по указанному номеру телефона нет заказов, пишем об этом,
		//иначе в цикле выводим данные по первым (в количестве perm.NmbrOfOutputOrders) заказам из слайса Data[], который содержит данные по заказам

		if len(User.Data) == 0 {
			output = "Заказы по указанным данным не найдены 🤷‍♂️\n"
		} else {
			var status string
			for i := 0; i < len(User.Data); i++ {
				//останавливаем цикл, если i достигло количества выводимых заказов
				if i == perm.NmbrOfOutputOrders {
					break
				}
				// каждой группе статуса соотносим словесное описание
				status = descriptionOfStatus(User.Data[i].Status.Group)
				// выбираем вариант заголовка в ответе бота перед статусом заказа/ов
				if i == 0 {
					output = titleOfStatus(len(User.Data), perm.NmbrOfOutputOrders)
					//switch len(User.Data) {
					//case 1:
					//	output = "Статус заказа ℹ\n"
					//case 2:
					//	output = "Найдено 2 заказа ℹ\n"
					//case 3:
					//	output = "Найдено 3 заказа ℹ\n"
					//default:
					//	output = "Найдено заказов: " + strconv.Itoa(len(User.Data)) + ". Покажу последние " + strconv.Itoa(perm.NmbrOfOutputOrders) + " ℹ\n"
					//}
				}
				//записываем в переменную output информацию по каждому найденному заказу, печатаем
				output = output +
					"\nНомер заказа: " + strconv.Itoa(User.Data[i].ID) +
					"\nБренд:  " + User.Data[i].Brand +
					"\nМодель: " + User.Data[i].Model +
					"\nЦена: " + fmt.Sprintf("%.2f", User.Data[i].Price) +
					"\nТекущий статус: " + status + "\n"
			}
		}
	}
	return output
}
