package models

import "time"

type LockCommit struct {
	Name string
	Time time.Time
	Pid  int
}
