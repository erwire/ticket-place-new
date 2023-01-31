package main

import (
	"fmt"
	"time"
)

func translateMonth(t time.Time) string {
	month := t.Month()
	switch month {
	case time.January:
		return "Января"
	case time.February:
		return "Февраля"
	case time.March:
		return "Марта"
	case time.April:
		return "Апреля"
	case time.May:
		return "Мая"
	case time.June:
		return "Июня"
	case time.July:
		return "Июля"
	case time.August:
		return "Августа"
	case time.September:
		return "Сентября"
	case time.October:
		return "Октября"
	case time.November:
		return "Ноября"
	case time.December:
		return "Декабря"
	default:
		return ""
	}

}
func getTime(timer time.Time) string {
	return timer.Format("15:04")
}
func getDate(date time.Time) string {
	return fmt.Sprint(date.Day()) + " " + translateMonth(date) + " " + fmt.Sprint(date.Year())
}

func GetDateTime(d string) (string, string) {
	t, err := time.Parse("2006-01-02 15:04:05", d)
	if err != nil {
		return "", ""
	}
	return getDate(t), getTime(t)
}
