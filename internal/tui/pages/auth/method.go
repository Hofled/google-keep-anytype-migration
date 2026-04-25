package auth

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

type MethodPage struct {
	*models.ModelInitOnce
	*models.PageIds

	focusedIndex int

	list list.Model

	appPageState state.AppPageStater
}

type item struct {
	title  string
	desc   string
	pageId models.PageId
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.title
}

func NewMethodPage(appPageState state.AppPageStater, apiKeyAuthPageId, challengeAuthPageId models.PageId) (*MethodPage, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	list := list.New([]list.Item{
		item{
			title: "API Key Authentication", desc: "Authenticate using API key", pageId: apiKeyAuthPageId,
		},
		item{
			title: "Challenge Authentication", desc: "Authenticate using temp challenge code", pageId: challengeAuthPageId,
		},
	}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Authentication Method"

	methodPage := &MethodPage{
		PageIds:      pageIds,
		list:         list,
		focusedIndex: 0,
		appPageState: appPageState,
	}

	methodPage.ModelInitOnce = models.NewModelInitOnce(methodPage)

	return methodPage, nil
}

func (m *MethodPage) Init() tea.Cmd {
	return nil
}

func (m *MethodPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch keyPress := msg.String(); keyPress {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				if cmd, err := m.appPageState.ShowPage(i.pageId); err != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MethodPage) View() tea.View {
	v := tea.NewView(m.list.View())
	return v
}
