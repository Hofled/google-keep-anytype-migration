package spaces

import (
	"context"
	"fmt"
	"strings"

	bubblesList "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/widgets/list"
)

type SpacesPageModel struct {
	*models.ModelInitOnce
	*models.PageIds

	authState   state.AppAuthStater
	windowState state.AppWindowStater

	spacesList *list.MultiSelectModel
}

type spaceListItem struct {
	space *rest.Space
}

func (slm *spaceListItem) FilterValue() string {
	return slm.space.Name
}

func (slm *spaceListItem) Title() string {
	var b strings.Builder

	if slm.space.Icon.Emoji != "" {
		b.WriteString(fmt.Sprintf("%s ", slm.space.Icon.Emoji))
	}

	b.WriteString(slm.space.Name)

	return b.String()
}

const noDescription = "<no description>"

func (slm *spaceListItem) Description() string {
	if len(slm.space.Description) == 0 {
		return noDescription
	} else {
		return slm.space.Description
	}
}

func NewSpacesModel(authState state.AppAuthStater, windowState state.AppWindowStater) (*SpacesPageModel, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	spacesPageModel := &SpacesPageModel{
		PageIds:     pageIds,
		authState:   authState,
		windowState: windowState,
	}

	spacesPageModel.ModelInitOnce = models.NewModelInitOnce(spacesPageModel)

	return spacesPageModel, nil

}

type spacesListMsg struct {
	list []*rest.Space
}

func (sm *SpacesPageModel) Init() tea.Cmd {
	return func() tea.Msg {
		client, err := rest.NewClient(sm.authState.GetAPIAddress())
		if err != nil || client == nil {
			return nil
		}

		client.SetApiKey(sm.authState.GetAPIKey())

		spacesList, err := client.ListSpaces(context.Background())
		if err != nil {
			return nil
		}

		return spacesListMsg{spacesList.Data}
	}
}

func (sm *SpacesPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgT := msg.(type) {
	case spacesListMsg:
		newSpacesList, err := constructSpacesList(msgT, sm.windowState.GetWindowWidth(), sm.windowState.GetWindowHeight())
		if err != nil {
			return sm, nil
		}

		sm.spacesList = newSpacesList
		return sm, nil
	}

	sm.spacesList, cmd = sm.spacesList.Update(msg)

	return sm, cmd
}

const spacesListTitle = "Choose Spaces For Import:"

func constructSpacesList(msg spacesListMsg, w, h int) (*list.MultiSelectModel, error) {
	spaces := make([]bubblesList.DefaultItem, len(msg.list))

	for i, space := range msg.list {
		spaces[i] = &spaceListItem{space}
	}

	spacesMultiSelect, err := list.NewMultiSelect(spaces, w, h)
	if err != nil {
		return nil, err
	}

	spacesMultiSelect.Title = spacesListTitle
	spacesMultiSelect.DisableQuitKeybindings()
	return spacesMultiSelect, nil
}

const loadingSpacesText = "Loading spaces..."

func (sm *SpacesPageModel) View() tea.View {
	var b strings.Builder

	if sm.spacesList == nil {
		b.WriteString(loadingSpacesText)
	} else {
		b.WriteString(sm.spacesList.View())
	}

	return tea.NewView(b.String())
}
