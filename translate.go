package calendar

import "time"

// Returns the name of the month in the selected language
func (cal *Calendar) getMonthDisplayName(month time.Month) string {
	if cal.opt.Language == RussianLangAbbr {
		return RussianMonths[month]
	}
	return month.String()
}

// Returns the array of the weekdays names in the selected language
func (cal *Calendar) getWeekdaysDisplayArray() [AmountOfDaysInWeek]string {
	if cal.opt.Language == RussianLangAbbr {
		return RussianWeekdaysAbbrs
	}
	return EnglishWeekdaysAbbrs
}
