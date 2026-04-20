package spaces

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	bubblesList "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/widgets/list"
)

type spacesListKeyMap struct {
	toggleSelection key.Binding
	confirmSpaces   key.Binding
	selectAll       key.Binding
	deselectAll     key.Binding
}

func newSpacesListKeyMap() *spacesListKeyMap {
	return &spacesListKeyMap{
		toggleSelection: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("␣/space", "toggle selection"),
		),
		confirmSpaces: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵/enter", "confirm selected"),
		),
		selectAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "select all"),
		),
		deselectAll: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("⇧+a", "deselect all"),
		),
	}
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

type SpacesPageModel struct {
	*models.ModelInitOnce
	*models.PageIds

	authState   state.AppAuthStater
	pageState   state.AppPageStater
	windowState state.AppWindowStater

	keyMap *spacesListKeyMap

	spacesList *list.MultiSelectModel
}

func NewSpacesModel(authState state.AppAuthStater, pageState state.AppPageStater, windowState state.AppWindowStater) (*SpacesPageModel, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	spacesPageModel := &SpacesPageModel{
		PageIds:     pageIds,
		authState:   authState,
		pageState:   pageState,
		windowState: windowState,
		keyMap:      newSpacesListKeyMap(),
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
	cmds := make([]tea.Cmd, 0)

	switch msgT := msg.(type) {
	case spacesListMsg:
		newSpacesList, err := constructSpacesList(msgT, sm.windowState.GetWindowWidth(), sm.windowState.GetWindowHeight(), sm.keyMap)
		if err != nil {
			return sm, nil
		}

		sm.spacesList = newSpacesList
		return sm, nil
	case tea.KeyPressMsg:
		if sm.spacesList.FilterState() != bubblesList.Filtering {
			k := msgT.Key()
			if key.Matches(k, sm.keyMap.confirmSpaces) {
				if cmd, err := sm.pageState.NextPage(); err != nil {
					cmds = append(cmds, cmd)
				}
			} else if key.Matches(k, sm.keyMap.selectAll) {
				sm.spacesList.SetAll(true)
			} else if key.Matches(k, sm.keyMap.deselectAll) {
				sm.spacesList.SetAll(false)
			}
		}
	}

	var cmd tea.Cmd
	sm.spacesList, cmd = sm.spacesList.Update(msg)
	cmds = append(cmds, cmd)

	return sm, tea.Batch(cmds...)
}

const spacesListTitle = "Choose Spaces For Import:"

func constructSpacesList(msg spacesListMsg, w, h int, keyMap *spacesListKeyMap) (*list.MultiSelectModel, error) {
	spaces := make([]bubblesList.DefaultItem, len(msg.list))

	for i, space := range msg.list {
		spaces[i] = &spaceListItem{space}
	}

	spacesMultiSelect, err := list.NewMultiSelect(spaces, w, h, keyMap.toggleSelection)
	if err != nil {
		return nil, err
	}

	keyBindings := []key.Binding{
		keyMap.toggleSelection,
		keyMap.confirmSpaces,
		keyMap.selectAll,
		keyMap.deselectAll,
	}

	spacesMultiSelect.Title = spacesListTitle
	spacesMultiSelect.DisableQuitKeybindings()
	spacesMultiSelect.AdditionalShortHelpKeys = func() []key.Binding {
		return keyBindings
	}
	spacesMultiSelect.AdditionalFullHelpKeys = func() []key.Binding {
		return keyBindings
	}

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
