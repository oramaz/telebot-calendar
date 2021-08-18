#  Overview
**Telebot-calendar** is an extension of the calendar inline-keyboard for Golang [Telebot.v3](http://github.com/tucnak/telebot "Telebot.v3")

```bash
go get github.com/oramaz/telebot-calendar
```



## Demo
![](https://imgur.com/mXRSytC.gif)

# Usage
Look at the initialization example
```go
import (
	tb_cal "github.com/oramaz/telebot-calendar"
	tb "gopkg.in/tucnak/telebot.v3"
)

b.Handle("/calendar", func(c tb.Context) error {
	calendar := tb_cal.NewCalendar(b, tb_cal.Options{})
	
	return c.Send("Select a date", &tb.ReplyMarkup{
		InlineKeyboard: calendar.GetKeyboard(),
	})
})
```
### Processing the input
You can get the selected date from the user like on the code bellow
```go
b.Handle(tb.OnText, func(c tb.Context) error {
	date, err := time.Parse("02.01.2006", c.Data())
	if err != nil {
		return err
	}
	
	// do smth with the received date
	fmt.Println(date)
	
	return nil
})
```
Pay attention to the special golang time layout parameter `"02.01.2006"` in the parsing of the received date ([More info](https://yourbasic.org/golang/format-parse-string-time-date-example/ "More info")).

[Full workable example code](https://pastebin.ubuntu.com/p/RfgzVbJ8sR/ "Full workable example code")

# Options
|  Name |Type  |  Defaut value | Description  |
| :------------ | :------------ | :------------ | :------------ |
|  InitialYear |  int |  current year | Initially active year |
|  InitialMonth|  time.Month |  current month | Initially active month |
|  YearRange |  [2]int | `[2]int{1970, 292277026596}` (time.Unix limits)  | The range of displayed years |
|  Language | string  |  `""` | The language of all designations. The value `"ru"` corresponds Russian, otherwise - English|
```go
calendar := tb_cal.NewCalendar(b, tb_cal.Options{
	YearRange: [2]int{2020, 2022},
	Language: "ru",
	. . .
})
```
# 
Have a good coding experience! 
I'll be happy to answer any questions and suggestions.
