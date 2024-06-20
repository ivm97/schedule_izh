package bot

import (
	"fmt"
	"strconv"
	"time"

	tgb "github.com/go-telegram-bot-api/telegram-bot-api"
)

const BTN_PREV = "<"
const BTN_NEXT = ">"

func GenerateCalendar(year int, month time.Month) tgb.InlineKeyboardMarkup {
	keyboard := tgb.InlineKeyboardMarkup{}
	keyboard = addMonthYearRow(year, month, keyboard)
	keyboard = addDaysNamesRow(keyboard)
	keyboard = generateMonth(year, int(month), keyboard)
	keyboard = addSpecialButtons(keyboard)
	return keyboard
}

func HandlerPrevButton(year int, month time.Month) (tgb.InlineKeyboardMarkup, int, time.Month) {
	if month != 1 {
		month--
	} else {
		month = 12
		year--
	}
	return GenerateCalendar(year, month), year, month
}

func HandlerNextButton(year int, month time.Month) (tgb.InlineKeyboardMarkup, int, time.Month) {
	if month != 12 {
		month++
	} else {
		year++
	}
	return GenerateCalendar(year, month), year, month
}

func addMonthYearRow(year int, month time.Month, keyboard tgb.InlineKeyboardMarkup) tgb.InlineKeyboardMarkup {
	var row []tgb.InlineKeyboardButton
	btn := tgb.NewInlineKeyboardButtonData(fmt.Sprintf("%s %v", month, year), "1")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func addDaysNamesRow(keyboard tgb.InlineKeyboardMarkup) tgb.InlineKeyboardMarkup {
	days := [7]string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var rowDays []tgb.InlineKeyboardButton
	for _, day := range days {
		btn := tgb.NewInlineKeyboardButtonData(day, day)
		rowDays = append(rowDays, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}

func generateMonth(year int, month int, keyboard tgb.InlineKeyboardMarkup) tgb.InlineKeyboardMarkup {
	firstDay := date(year, month, 0)
	amountDaysInMonth := date(year, month+1, 0).Day()

	weekday := int(firstDay.Weekday())
	rowDays := []tgb.InlineKeyboardButton{}
	for i := 1; i <= weekday; i++ {
		btn := tgb.NewInlineKeyboardButtonData(" ", string(i))
		rowDays = append(rowDays, btn)
	}

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
			amountWeek = 0
			rowDays = []tgb.InlineKeyboardButton{}
		}

		day := strconv.Itoa(i)
		if len(day) == 1 {
			day = fmt.Sprintf("0%v", day)
		}
		monthStr := strconv.Itoa(month)
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%v", monthStr)
		}

		btnText := fmt.Sprintf("%v", i)
		if time.Now().Day() == i {
			btnText = fmt.Sprintf("%v!", i)
		}
		btn := tgb.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("%v.%v.%v", year, monthStr, day))
		rowDays = append(rowDays, btn)
		amountWeek++
	}
	for i := 1; i <= 7-amountWeek; i++ {
		btn := tgb.NewInlineKeyboardButtonData(" ", string(i))
		rowDays = append(rowDays, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)

	return keyboard
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func addSpecialButtons(keyboard tgb.InlineKeyboardMarkup) tgb.InlineKeyboardMarkup {
	var rowDays = []tgb.InlineKeyboardButton{}
	btnPrev := tgb.NewInlineKeyboardButtonData(BTN_PREV, BTN_PREV)
	btnNext := tgb.NewInlineKeyboardButtonData(BTN_NEXT, BTN_NEXT)
	rowDays = append(rowDays, btnPrev, btnNext)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}
