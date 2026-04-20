package spaces

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

type SpacesPageModel struct {
	*models.ModelInitOnce
	*models.PageIds

	authState state.AppAuthStater

	spacesList *list.Model
}

type spaceListItem struct {
	*rest.Space
}

func (slm *spaceListItem) FilterValue() string {
	return slm.Name
}

func NewSpacesModel(authState state.AppAuthStater) (*SpacesPageModel, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	spacesPageModel := &SpacesPageModel{
		PageIds:   pageIds,
		authState: authState,
	}

	spacesPageModel.ModelInitOnce = models.NewModelInitOnce(spacesPageModel)

	return spacesPageModel, nil

}

type spacesListMsg struct {
	list []*rest.Space
}

func (sm *SpacesPageModel) Init() tea.Cmd {
	return func() tea.Msg {
		client := sm.authState.GetClient()
		if client == nil {
			return nil
		}

		spacesList, err := client.ListSpaces(context.Background())
		if err != nil {
			return nil
		}

		return spacesListMsg{spacesList.Data}
	}
}

func (sm *SpacesPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgT := msg.(type) {
	case spacesListMsg:
		sm.spacesList = constructSpacesList(msgT)
	}

	return sm, nil
}

func constructSpacesList(msg spacesListMsg) *list.Model {
	spaces := make([]list.Item, len(msg.list))

	for i, space := range msg.list {
		spaces[i] = &spaceListItem{space}
	}

	spacesList := list.New(spaces, list.NewDefaultDelegate(), 0, 0)
	return &spacesList
}

func (sm *SpacesPageModel) View() tea.View {
	var b strings.Builder

	return tea.NewView(b.String())
}
