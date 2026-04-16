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

func NewMethodPage(appPageState state.AppPageStater, apiKeyAuthPageId models.PageId) (*MethodPage, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	list := list.New([]list.Item{item{title: "API Key Authentication", desc: "Authenticate using API key", pageId: apiKeyAuthPageId}}, list.NewDefaultDelegate(), 0, 0)
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

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *MethodPage) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *MethodPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch keyPress := msg.String(); keyPress {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.appPageState.ShowPage(i.pageId)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (m *MethodPage) View() tea.View {
	v := tea.NewView(m.list.View())
	return v
}
