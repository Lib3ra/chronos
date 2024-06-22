package model

import "time"

type DayEntry struct {
	Date  time.Time
	Start time.Time
	End   time.Time
}
