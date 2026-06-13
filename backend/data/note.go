package data

import "time"

type NoteStatus string

const (
	NSActive   NoteStatus = "Active"
	NSArchived NoteStatus = "Archived"
)

type Note struct {
	Id            int
	Title         string
	Body          string
	Status        NoteStatus
	CreatedTs     time.Time
	LastUpdatedTs time.Time
}
