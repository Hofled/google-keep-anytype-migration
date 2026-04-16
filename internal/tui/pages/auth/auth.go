package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"
	"github.com/google/uuid"
)

var (
	disabledTextStyle = lipgloss.NewStyle().Strikethrough(true).Faint(true)
	authErrStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type AuthPage struct {
	*models.ModelInitOnce

	id         models.PageId
	nextPageId *models.PageId
	prevPageId *models.PageId

	addrInput    textinput.Model
	keyInput     textinput.Model
	errorMsg     string
	connected    atomic.Bool
	focusedIndex int

	setClientOnce sync.Once

	appAuthState state.AppAuthStater
	appViewState state.AppPageStater
}

type authResultMsg struct {
	success bool
	err     error
}

func NewAuthPage(appAuthState state.AppAuthStater, appViewState state.AppPageStater, prevPageId, nextPageId *models.PageId) (*AuthPage, error) {
	addrInput := textinput.New()
	addrInput.SetValue("http://localhost:31009")
	addrInput.Placeholder = "http://localhost:31009"
	addrInput.Focus()
	addrInput.SetWidth(50)

	keyInput := textinput.New()
	keyInput.Placeholder = "Your API Key"
	keyInput.SetWidth(50)

	id, err := models.NewPageID()
	if err != nil {
		return nil, fmt.Errorf("failed creating auth page: %w", err)
	}

	authPage := &AuthPage{
		id:            id,
		nextPageId:    nextPageId,
		prevPageId:    prevPageId,
		addrInput:     addrInput,
		keyInput:      keyInput,
		focusedIndex:  0,
		appAuthState:  appAuthState,
		appViewState:  appViewState,
		setClientOnce: sync.Once{},
	}

	authPage.ModelInitOnce = models.NewModelInitOnce(authPage)

	return authPage, nil
}

func (a *AuthPage) Init() tea.Cmd {
	return textinput.Blink
}

func (a *AuthPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			a.handleNavigation(msg.String())
			return a, nil
		case "enter":
			if a.focusedIndex == 2 {
				return a, a.connect()
			} else if a.focusedIndex == 3 && a.CanProceed() {
				a.appViewState.NextPage()
				// Proceed to next page (handled by app)
				return a, nil
			}
		}
	case authResultMsg:
		if msg.success {
			a.errorMsg = ""
		} else {
			a.errorMsg = msg.err.Error()
		}
		return a, nil
	}

	// Update inputs
	switch a.focusedIndex {
	case 0:
		a.addrInput, cmd = a.addrInput.Update(msg)
	case 1:
		a.keyInput, cmd = a.keyInput.Update(msg)
	}

	return a, cmd
}

func (a *AuthPage) View() tea.View {
	var b strings.Builder

	b.WriteString("Authentication\n\n")

	b.WriteString(fmt.Sprintf("API Address: %s\n", a.addrInput.View()))
	b.WriteString(fmt.Sprintf("API Key: %s\n\n", a.keyInput.View()))

	connectLabel := "Connect" // TODO add spinner during connection
	if a.connected.Load() {
		connectLabel = disabledTextStyle.Render(connectLabel)
	}
	if a.focusedIndex == 2 {
		connectLabel = "[" + connectLabel + "]"
	}
	b.WriteString(fmt.Sprintf("%s\n", connectLabel))

	nextLabel := "Next"
	if !a.CanProceed() {
		nextLabel = disabledTextStyle.Render(nextLabel)
	}
	if a.focusedIndex == 3 {
		nextLabel = "[" + nextLabel + "]"
	}
	b.WriteString(fmt.Sprintf("%s\n\n", nextLabel))

	if a.errorMsg != "" {
		b.WriteString(authErrStyle.Render("❌ Error: "+a.errorMsg) + "\n")
	}

	if a.connected.Load() {
		b.WriteString("✓ Connected successfully!\n")
	}

	return tea.NewView(b.String())
}

func (a *AuthPage) CanProceed() bool {
	return a.connected.Load()
}

func (a *AuthPage) GetData() any {
	return map[string]string{
		"addr": a.addrInput.Value(),
		"key":  a.keyInput.Value(),
	}
}

func (a *AuthPage) NextPageId() models.PageId {
	if a.nextPageId != nil {
		return *a.nextPageId
	} else {
		return models.PageId(uuid.Nil)
	}
}

func (a *AuthPage) PrevPageId() models.PageId {
	if a.prevPageId != nil {
		return *a.prevPageId
	} else {
		return models.PageId(uuid.Nil)
	}
}

func (a *AuthPage) ID() models.PageId {
	return a.id
}

func (a *AuthPage) handleNavigation(key string) {
	switch key {
	case "tab":
		a.focusedIndex = (a.focusedIndex + 1) % 4
	case "shift+tab":
		a.focusedIndex = (a.focusedIndex - 1 + 4) % 4
	case "up":
		a.focusedIndex = (a.focusedIndex - 1 + 4) % 4
	case "down":
		a.focusedIndex = (a.focusedIndex + 1) % 4
	}

	// Update focus
	a.addrInput.Blur()
	a.keyInput.Blur()

	switch a.focusedIndex {
	case 0:
		a.addrInput.Focus()
	case 1:
		a.keyInput.Focus()
	}
}

func (a *AuthPage) connect() tea.Cmd {
	if a.connected.Load() {
		return nil
	}

	return tea.Cmd(func() tea.Msg {
		authenticatedClient := anytype.NewClient(anytype.WithBaseURL(a.addrInput.Value()), anytype.WithAppKey(a.keyInput.Value()))
		if authenticatedClient == nil {
			return authResultMsg{success: false, err: fmt.Errorf("failed creating client")}
		}

		if _, err := authenticatedClient.Spaces().List(context.Background()); err != nil {
			return authResultMsg{success: false, err: fmt.Errorf("invalid API key: %w", err)}
		}

		a.setClientOnce.Do(func() {
			a.connected.Store(true)
			a.appAuthState.SetClient(authenticatedClient)
		})

		return authResultMsg{success: true, err: nil}
	})
}
