package clock

import (
	"fmt"
	"time"
)

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct {
	IsAsia bool
}

func (fc FixedClocker) Now() time.Time {
	// UTCで指定された日時
	utcTime := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)

	if fc.IsAsia {
		// Asia / Tokyoのタイムゾーンに変換
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			fmt.Println("Error:", err)
			return time.Time{}
		}
		tokyoTime := utcTime.In(loc)
		return tokyoTime
	}

	return utcTime
}
