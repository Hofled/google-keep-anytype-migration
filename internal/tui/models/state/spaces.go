package state

import "github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"

type ImportSpacesStater interface {
	SetSelectedSpaces(spaces []*rest.Space)
	SelectedSpaces() []*rest.Space
}

type ImportSpacesState struct {
	selectedSpaces []*rest.Space
}

func (iss *ImportSpacesState) SetSelectedSpaces(spaces []*rest.Space) {
	iss.selectedSpaces = spaces
}

func (iss *ImportSpacesState) SelectedSpaces() []*rest.Space {
	return iss.selectedSpaces
}
