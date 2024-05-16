package models

import "time"

type StdLine struct {
	Cwd  string
	Time time.Time
	Data string
}
