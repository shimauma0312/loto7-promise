package result

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

	// 2025年10月18日を第648回の基準日
	baseFriday := time.Date(2025, 10, 18, 0, 0, 0, 0, time.Local)
	baseNumber := 648

	targetFriday := getLastFriday(today)

	weeksDiff := int(targetFriday.Sub(baseFriday).Hours() / (24 * 7))

	return baseNumber + weeksDiff
}
