package bot

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivm97/schedule_izh/models"
)

//TODO
//Сделать нормальнo

func (app *application) start() {

	bot, err := tgb.NewBotAPI(app.settings.Bot.Token)
	if err != nil {
		app.eLog.Fatal(err)
	}

	bot.Debug = true
	app.iLog.Println("Success. Bot %s started!", bot.Self.UserName)
	u := tgb.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			app.iLog.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.CallbackQuery != nil {
				go func(replyTo int, data string, chatId int64) {
					ses, ok := app.sess.Get(chatId)
					if !ok {
						app.eLog.Println("Ошибка, пользователь дошел до инлайн взаимодействия без сессии")
					} else {
						switch {
						case ses.DepDay == "":
							if ok = dateCheck(data); ok {
								ses.DepDay = data
								app.sess.Set(chatId, ses, time.Minute*8)
								msg := msgBuilder(chatId, "Теперь выберите дату возвращения!", replyTo, ses.Keyboard)
								bot.Send(msg)
							} else {
								msg := msgBuilder(chatId, "Не верный формат даты... ", replyTo, ses.Keyboard)
								bot.Send(msg)
							}
						case ses.ArrDay == "":
							if ok = dateCheck(data); ok {
								ses.ArrDay = data
								app.sess.Set(chatId, ses, time.Minute*8)
								msg := msgBuilder(chatId, fmt.Sprint(ses), replyTo, nil)
								bot.Send(msg)
							}
						}
					}
				}(update.Message.MessageID, update.CallbackQuery.Data, update.Message.Chat.ID)
			}
			go func(replyTo int, message string, chatId int64) {
				if message == "/отмена" {
					err := app.sess.Delete(chatId)
					if err != nil {
						msg := msgBuilder(chatId, "Сессия еще не была инициализирована!", replyTo, nil)
						bot.Send(msg)
					}
				}

				if message == "/что" {
					_, in := app.sess.Get(chatId)
					msg := msgBuilder(chatId, fmt.Sprint(in), replyTo, nil)
					bot.Send(msg)
				}

				ses, ok := app.sess.Get(update.Message.Chat.ID)
				if !ok && update.Message.Text == "/routes" {
					//Инициализируем сессию и передаем соответствующие указания пользователю
					app.sess.Set(chatId, models.SearchData{}, time.Minute*8)
					msg := msgBuilder(chatId, "Выберите точку отправки с помощью кнопок!", replyTo, cityKB())
					bot.Send(msg)
				} else if !ok {
					msg := msgBuilder(chatId, "Введите /routes для проложения маршрута!", replyTo, nil)
					bot.Send(msg)
				} else {

					//Логика проверки значений
					switch {
					case ses.Dep == "":
						if val, ok := cityHandler(message); ok {
							ses.Dep = val
							app.sess.Set(chatId, ses, time.Minute*8)
							msg := msgBuilder(chatId, "Теперь введите точку прибытия!", replyTo, cityKB())
							bot.Send(msg)
						} else {
							msg := msgBuilder(chatId, "В такой город не летаем!", replyTo, cityKB())
							bot.Send(msg)
						}
					case ses.Arr == "":
						if val, ok := cityHandler(message); ok && ses.Dep != val {
							ses.Arr = val
							app.sess.Set(chatId, ses, time.Minute*8)
							if ses.Keyboard == nil {
								ses.Keyboard = GenerateCalendar(time.Now().Year(), time.Now().Month())
							}
							msg := msgBuilder(chatId, "Осталось выбрать дату вылета! Используйте инлайн кнопки:", replyTo, ses.Keyboard)
							bot.Send(msg)
						} else if ses.Arr == ses.Dep {
							msg := msgBuilder(chatId, "В пределах одного города можно пройтись пешком!", replyTo, cityKB())
							bot.Send(msg)
						} else {
							msg := msgBuilder(chatId, "Неверно введена точка прибытия! Выберите с помощью кнопок!", replyTo, cityKB())
							bot.Send(msg)
						}

					case ses.DepDay == "":
						if ok := dateCheck(message); ok {
							ses.DepDay = message
							app.sess.Set(chatId, ses, time.Minute*8)
							msg := msgBuilder(chatId, "Теперь введите дату прибытия с помощью инлайн клавиатуры!", replyTo, ses.Keyboard)
							bot.Send(msg)
						}
					}

				}

			}(update.Message.MessageID, update.Message.Text, update.Message.Chat.ID)

		}
	}
}

func cityKB() tgb.ReplyKeyboardMarkup {
	var numericKeyboard = tgb.NewReplyKeyboard(
		tgb.NewKeyboardButtonRow(
			tgb.NewKeyboardButton("Ижевск"),
			tgb.NewKeyboardButton("Иркутск"),
			tgb.NewKeyboardButton("Калининград"),
			tgb.NewKeyboardButton("Когалым"),
			tgb.NewKeyboardButton("Махачкала"),
		),
		tgb.NewKeyboardButtonRow(
			tgb.NewKeyboardButton("Москва"),
			tgb.NewKeyboardButton("Домодедово, Москва"),
			tgb.NewKeyboardButton("Санкт-Петербург"),
			tgb.NewKeyboardButton("Пулково, Санкт-Петербург"),
			tgb.NewKeyboardButton("Сочи"),
		),
	)

	return numericKeyboard
}

func msgBuilder(chatId int64, message string, replyTo int, kb any) tgb.MessageConfig {
	msg := tgb.NewMessage(chatId, message)
	msg.ReplyToMessageID = replyTo
	if kb != nil {
		msg.ReplyMarkup = kb
	} else {
		msg.ReplyMarkup = tgb.NewRemoveKeyboard(true)
	}

	return msg
}

func cityHandler(message string) (string, bool) {
	cmap := map[string]string{
		"Ижевск":                   "IJK",
		"Иркутск":                  "IKT",
		"Калининград":              "KGD",
		"Когалым":                  "KGP",
		"Махачкала":                "MCX",
		"Москва":                   "MOW",
		"Домодедово, Москва":       "DME",
		"Санкт-Петербург":          "LED",
		"Пулково, Санкт-Петербург": "LED",
		"Сочи":                     "AER",
	}

	val, ok := cmap[message]
	return val, ok
}

func dateCheck(message string) bool {

	slice := strings.Split(message, ".")
	for _, e := range slice {
		for _, subE := range e {
			if !unicode.IsDigit(subE) {
				return false
			}
		}
	}
	return true
}
