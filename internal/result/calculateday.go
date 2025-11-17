package result

import (
	"time"
)

func getFriday(t time.Time) time.Time {
	offset := (5 - int(t.Weekday()) + 7) % 7
	return t.AddDate(0, 0, offset)
}

func getLastFriday(t time.Time) time.Time {
	// 今日が金曜日の場合は今日を返す
	if t.Weekday() == time.Friday {
		return t
	}

	// 直前の金曜日を計算
	// 日曜日=0, 月曜日=1, ..., 金曜日=5, 土曜日=6
	currentDay := int(t.Weekday())
	fridayDay := int(time.Friday)

	var daysBack int
	if currentDay > fridayDay {
		// 今日が土曜日（6）の場合: 6 - 5 = 1日前
		daysBack = currentDay - fridayDay
	} else {
		// 今日が日〜木曜日の場合
		// 例: 水曜日(3)の場合: 3 + 7 - 5 = 5日前
		daysBack = currentDay + 7 - fridayDay
	}

	return t.AddDate(0, 0, -daysBack)
}

func NewNumber() int {
	today := time.Now()

	// 2025年10月17日（金）を第648回の基準日
	baseFriday := time.Date(2025, 10, 17, 0, 0, 0, 0, time.Local)
	baseNumber := 648

	targetFriday := getLastFriday(today)

	weeksDiff := int(targetFriday.Sub(baseFriday).Hours() / (24 * 7))

	return baseNumber + weeksDiff
}

// GetDrawDate 抽選回数から実際の抽選日を計算
func GetDrawDate(drawNumber int) string {
	// 2025年10月17日（金）を第648回の基準日
	baseFriday := time.Date(2025, 10, 17, 0, 0, 0, 0, time.Local)
	baseNumber := 648

	// 基準日からの差分を計算
	weeksDiff := drawNumber - baseNumber

	// 抽選日を計算（毎週金曜日）
	drawDate := baseFriday.AddDate(0, 0, weeksDiff*7)

	return drawDate.Format("2006-01-02")
}
