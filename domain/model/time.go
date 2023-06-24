package model

import "time"

type Time struct {
	Time time.Time
}

func NewTime(time time.Time) *Time {
	return &Time{Time: time}
}

// yyyy/MM/dd hh:mm:ssの形にフォーマットする
func (t *Time) Format() string {
	const layout = "2006/01/02 15:04:05"
	return t.Time.Format(layout)
}
