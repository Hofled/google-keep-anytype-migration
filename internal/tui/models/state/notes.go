package state

import "github.com/Hofled/go-google-keep-anytype-migration/internal/googlekeep"

type NotesStater interface {
	ParsedNotes() []googlekeep.Note
	SetParsedNotes(notes []googlekeep.Note)
}

type NotesState struct {
	parsedNotes []googlekeep.Note
}

func (ns *NotesState) ParsedNotes() []googlekeep.Note {
	return ns.parsedNotes
}

func (ns *NotesState) SetParsedNotes(notes []googlekeep.Note) {
	ns.parsedNotes = notes
}
