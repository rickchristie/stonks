package api

import (
	"time"

	"stonks/data"
)

type Note struct {
	Id            int             `json:"id"`
	Title         string          `json:"title"`
	Body          string          `json:"body"`
	Status        data.NoteStatus `json:"status"`
	CreatedTs     time.Time       `json:"createdTs"`
	LastUpdatedTs time.Time       `json:"lastUpdatedTs"`
}

func NoteFromData(note *data.Note) *Note {
	if note == nil {
		return nil
	}
	return &Note{
		Id:            note.Id,
		Title:         note.Title,
		Body:          note.Body,
		Status:        note.Status,
		CreatedTs:     note.CreatedTs,
		LastUpdatedTs: note.LastUpdatedTs,
	}
}

func NotesFromData(notes []*data.Note) []*Note {
	ret := make([]*Note, 0, len(notes))
	for _, note := range notes {
		ret = append(ret, NoteFromData(note))
	}
	return ret
}
