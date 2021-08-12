package calendar

import "time"

var (
	// RussianMonths is a map of months names translated to Russian
	RussianMonths map[time.Month]string = map[time.Month]string{
		time.January:   "Январь",
		time.February:  "Февраль",
		time.March:     "Март",
		time.April:     "Апрель",
		time.May:       "Май",
		time.June:      "Июнь",
		time.July:      "Июль",
		time.August:    "Август",
		time.September: "Сентябрь",
		time.October:   "Октябрь",
		time.November:  "Ноябрь",
		time.December:  "Декабрь",
	}

	// EnglishWeekdaysAbbrs is an array of weekdays names
	EnglishWeekdaysAbbrs = [AmountOfDaysInWeek]string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}

	// RussianWeekdaysAbbrs is an array of weekdays names translated to Russian
	RussianWeekdaysAbbrs = [AmountOfDaysInWeek]string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
)

// Unix years range ([1970, 292277026596])
var (
	MinYearLimit = time.Unix(0, 0).Year()
	MaxYearLimit = time.Unix(1<<63-1, 0).Year()
)

// AmountOfDaysInWeek is a constant that represents an amount of days in week
const AmountOfDaysInWeek = 7

// RussianLangAbbr is a constant that represents the Russian language abbreviation
// that could be passed in Options.Language
const RussianLangAbbr = "ru"
