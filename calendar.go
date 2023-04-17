package calendar

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tb "gopkg.in/telebot.v3"
)

// NewCalendar builds and returns a Calendar
func NewCalendar(b *tb.Bot, opt Options) *Calendar {
	if opt.YearRange == [2]int{0, 0} {
		opt.YearRange = [2]int{MinYearLimit, MaxYearLimit}
	}
	if opt.InitialYear == 0 {
		opt.InitialYear = time.Now().Year()
	}
	if opt.InitialMonth == 0 {
		opt.InitialMonth = time.Now().Month()
	}

	err := opt.validate()
	if err != nil {
		panic(err)
	}

	return &Calendar{
		Bot:       b,
		kb:        make([][]tb.InlineButton, 0),
		opt:       &opt,
		currYear:  opt.InitialYear,
		currMonth: opt.InitialMonth,
	}
}

// Calendar represents the main object
type Calendar struct {
	Bot       *tb.Bot
	opt       *Options
	kb        [][]tb.InlineButton
	currYear  int
	currMonth time.Month
}

// Options represents a struct for passing optional
// properties for customizing a calendar keyboard
type Options struct {
	// The year that will be initially active in the calendar.
	// Default value - today's year
	InitialYear int

	// The month that will be initially active in the calendar
	// Default value - today's month
	InitialMonth time.Month

	// The range of displayed years
	// Default value - {1970, 292277026596} (time.Unix years range)
	YearRange [2]int

	// The language of all designations.
	// If equals "ru" the designations would be Russian,
	// otherwise - English
	Language string
}

// GetKeyboard builds the calendar inline-keyboard
func (cal *Calendar) GetKeyboard() [][]tb.InlineButton {
	cal.clearKeyboard()

	cal.addMonthYearRow()
	cal.addWeekdaysRow()
	cal.addDaysRows()
	cal.addControlButtonsRow()

	return cal.kb
}

// Clears the calendar's keyboard
func (cal *Calendar) clearKeyboard() {
	cal.kb = make([][]tb.InlineButton, 0)
}

// Builds a full row width button with a displayed month's name
// The button represents a list of all months when clicked
func (cal *Calendar) addMonthYearRow() {
	var row []tb.InlineButton

	btn := tb.InlineButton{
		Unique: genUniqueParam("month_year_btn"),
		Text:   fmt.Sprintf("%s %v", cal.getMonthDisplayName(cal.currMonth), cal.currYear),
	}

	cal.Bot.Handle(&btn, func(ctx tb.Context) error {
		cal.Bot.Edit(ctx.Message(), ctx.Message().Text, &tb.ReplyMarkup{
			InlineKeyboard: cal.getMonthPickKeyboard(),
		})

		ctx.Respond()
		return nil
	})

	row = append(row, btn)
	cal.addRowToKeyboard(&row)
}

// Builds a keyboard with a list of months to pick
func (cal *Calendar) getMonthPickKeyboard() [][]tb.InlineButton {
	cal.clearKeyboard()

	var row []tb.InlineButton

	// Generating a list of months
	for i := 1; i <= 12; i++ {
		monthName := cal.getMonthDisplayName(time.Month(i))
		monthBtn := tb.InlineButton{
			Unique: genUniqueParam("month_pick_" + fmt.Sprint(i)),
			Text:   monthName, Data: strconv.Itoa(i),
		}

		cal.Bot.Handle(&monthBtn, func(ctx tb.Context) error {
			monthNum, err := strconv.Atoi(ctx.Data())
			if err != nil {
				log.Fatal(err)
			}
			cal.currMonth = time.Month(monthNum)

			// Show the calendar keyboard with the active selected month back
			cal.Bot.Edit(ctx.Message(), ctx.Message().Text, &tb.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			})

			ctx.Respond()
			return nil
		})

		row = append(row, monthBtn)

		// Arranging the months in 2 columns
		if i%2 == 0 {
			cal.addRowToKeyboard(&row)
			row = []tb.InlineButton{} // empty row
		}
	}

	return cal.kb
}

// Builds a row of non-clickable buttons
// that display weekdays names
func (cal *Calendar) addWeekdaysRow() {
	var row []tb.InlineButton

	for i, wd := range cal.getWeekdaysDisplayArray() {
		btn := tb.InlineButton{Unique: genUniqueParam("weekday_" + fmt.Sprint(i)), Text: wd}
		cal.Bot.Handle(&btn, ignoreQuery)
		row = append(row, btn)
	}

	cal.addRowToKeyboard(&row)
}

// Builds a table of clickable cells (buttons) - active month's days
func (cal *Calendar) addDaysRows() {
	beginningOfMonth := time.Date(cal.currYear, cal.currMonth, 1, 0, 0, 0, 0, time.UTC)
	amountOfDaysInMonth := beginningOfMonth.AddDate(0, 1, -1).Day()

	var row []tb.InlineButton

	// Calculating the number of empty buttons that need to be inserted forward
	weekdayNumber := int(beginningOfMonth.Weekday())
	if weekdayNumber == 0 && cal.opt.Language == RussianLangAbbr { // russian Sunday exception
		weekdayNumber = 7
	}

	// The difference between English and Russian weekdays order
	// en: Sunday (0), Monday (1), Tuesday (3), ...
	// ru: Monday (1), Tuesday (2), ..., Sunday (7)
	if cal.opt.Language != RussianLangAbbr {
		weekdayNumber++
	}

	// Inserting empty buttons forward
	for i := 1; i < weekdayNumber; i++ {
		cal.addEmptyCell(&row)
	}

	// Inserting month's days' buttons
	for i := 1; i <= amountOfDaysInMonth; i++ {
		dayText := strconv.Itoa(i)
		cell := tb.InlineButton{
			Unique: genUniqueParam("day_" + fmt.Sprint(i)),
			Text:   dayText, Data: dayText,
		}

		cal.Bot.Handle(&cell, func(ctx tb.Context) error {
			dayInt, err := strconv.Atoi(ctx.Data())
			if err != nil {
				return err
			}
			ctx.Message().Payload = cal.genDateStrFromDay(dayInt)

			upd := tb.Update{Message: ctx.Message()}
			cal.Bot.ProcessUpdate(upd)

			ctx.Respond()
			return nil
		})

		row = append(row, cell)

		if len(row)%AmountOfDaysInWeek == 0 {
			cal.addRowToKeyboard(&row)
			row = []tb.InlineButton{} // empty row
		}
	}

	// Inseting empty buttons at the end
	if len(row) > 0 {
		for i := len(row); i < AmountOfDaysInWeek; i++ {
			cal.addEmptyCell(&row)
		}
		cal.addRowToKeyboard(&row)
	}
}

// Builds a row of  control buttons for swiping the calendar
func (cal *Calendar) addControlButtonsRow() {
	var row []tb.InlineButton

	prev := tb.InlineButton{Unique: genUniqueParam("prev_month"), Text: "＜"}

	// Hide "prev" button if it rests on the range
	if cal.currYear <= cal.opt.YearRange[0] && cal.currMonth == 1 {
		prev.Text = ""
	} else {
		cal.Bot.Handle(&prev, func(ctx tb.Context) error {
			// Additional protection against entering the years ranges
			if cal.currMonth > 1 {
				cal.currMonth--
			} else {
				cal.currMonth = 12
				if cal.currYear > cal.opt.YearRange[0] {
					cal.currYear--
				}
			}

			cal.Bot.Edit(ctx.Message(), ctx.Message().Text, &tb.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			})

			ctx.Respond()
			return nil
		})
	}

	next := tb.InlineButton{Unique: genUniqueParam("next_month"), Text: "＞"}

	// Hide "next" button if it rests on the range
	if cal.currYear >= cal.opt.YearRange[1] && cal.currMonth == 12 {
		next.Text = ""
	} else {
		cal.Bot.Handle(&next, func(ctx tb.Context) error {
			// Additional protection against entering the years ranges
			if cal.currMonth < 12 {
				cal.currMonth++
			} else {
				if cal.currYear < cal.opt.YearRange[1] {
					cal.currYear++
				}
				cal.currMonth = 1
			}

			cal.Bot.Edit(ctx.Message(), ctx.Message().Text, &tb.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			})

			ctx.Respond()
			return nil
		})
	}

	row = append(row, prev, next)
	cal.addRowToKeyboard(&row)
}

// Returns a formatted date string from the selected date
func (cal *Calendar) genDateStrFromDay(day int) string {
	return time.Date(cal.currYear, cal.currMonth, day,
		0, 0, 0, 0, time.UTC).Format("02.01.2006")
}

// Generates a InlineButton.Unique param via concatenating
// base name with a random string.
// It is necessary to exclude the overlap of links between buttons
func genUniqueParam(base string) string {
	return fmt.Sprintf("%s_%s", base, randSequence(8))
}

// Utility function for passing a row to the calendar's keyboard
func (cal *Calendar) addRowToKeyboard(row *[]tb.InlineButton) {
	cal.kb = append(cal.kb, *row)
}

// Inserts an empty button that doesn't process queires
// into the keyboard row
func (cal *Calendar) addEmptyCell(row *[]tb.InlineButton) {
	cell := tb.InlineButton{Unique: genUniqueParam("empty_cell"), Text: " "}
	cal.Bot.Handle(&cell, ignoreQuery)
	*row = append(*row, cell)
}

// Query stub
func ignoreQuery(ctx tb.Context) error {
	ctx.Respond()
	return nil
}
