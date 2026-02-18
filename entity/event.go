package entity

import "time"

type Event struct {
	DocumentID  string
	Type        string
	FileName    string
	Timestamp   time.Time
	ContentType *string
}
