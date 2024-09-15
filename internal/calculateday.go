package internal

import (
	"time"
)

func getFriday(t time.Time) time.Time {
	offset := (5 - int(t.Weekday()) + 7) % 7
	return t.AddDate(0, 0, offset)
}

func getLastFriday(t time.Time) time.Time {
	if t.Weekday() == time.Friday {
		return t
	}
	return getFriday(t.AddDate(0, 0, -7))
}

func NewNumber() int {
	today := time.Now()

	baseFriday := time.Date(2024, 9, 6, 0, 0, 0, 0, time.Local)
	baseNumber := 591

	targetFriday := getLastFriday(today)

	weeksDiff := int(targetFriday.Sub(baseFriday).Hours() / (24 * 7))

	return baseNumber + weeksDiff
}
