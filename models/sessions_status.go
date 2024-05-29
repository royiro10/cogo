package models

import "time"

type SessionStatus struct {
	SessionStatus        string
	LastCommand          string
	LastActionTime       time.Time
	ExecutedCommandCount uint
	ExecuteQueueSize     uint
	OutputViewSize       uint
}
