package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/styles"
)

const (
	addrInputFocusIndex = iota
	keyInputFocusIndex
	connectButtFocusIndex
	nextButtonFocusIndex
	prevButtonFocusIndex
)

type ApiKeyAuthPage struct {
	*models.ModelInitOnce
	*models.PageIds

	addrInput    textinput.Model
	keyInput     textinput.Model
	errorMsg     string
	connected    atomic.Bool
	focusedIndex int

	setClientOnce sync.Once

	appAuthState state.AppAuthStater
	appPageState state.AppPageStater
}

type authResultMsg struct {
	success bool
	err     error
}

func NewApiKeyAuthPage(appAuthState state.AppAuthStater, appPageState state.AppPageStater) (*ApiKeyAuthPage, error) {
	addrInput := textinput.New()
	addrInput.SetValue("http://localhost:31009")
	addrInput.Placeholder = "http://localhost:31009"
	addrInput.Focus()
	addrInput.SetWidth(50)

	keyInput := textinput.New()
	keyInput.Placeholder = "Your API Key"
	keyInput.SetWidth(50)

	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	authPage := &ApiKeyAuthPage{
		PageIds:       pageIds,
		addrInput:     addrInput,
		keyInput:      keyInput,
		focusedIndex:  0,
		appAuthState:  appAuthState,
		appPageState:  appPageState,
		setClientOnce: sync.Once{},
	}

	authPage.ModelInitOnce = models.NewModelInitOnce(authPage)

	return authPage, nil
}

func (a *ApiKeyAuthPage) Init() tea.Cmd {
	return textinput.Blink
}

func (a *ApiKeyAuthPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if cmd, err := a.appPageState.NextPage(); err == nil {
					return a, cmd
				}

				return a, nil
			} else if a.focusedIndex == 4 {
				if cmd, err := a.appPageState.PrevPage(); err == nil {
					return a, cmd
				}

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

	switch a.focusedIndex {
	case addrInputFocusIndex:
		a.addrInput, cmd = a.addrInput.Update(msg)
	case keyInputFocusIndex:
		a.keyInput, cmd = a.keyInput.Update(msg)
	}

	return a, cmd
}

func (a *ApiKeyAuthPage) View() tea.View {
	var b strings.Builder

	b.WriteString("API Key Authentication\n\n")

	fmt.Fprintf(&b, "API Address: %s\n", a.addrInput.View())
	fmt.Fprintf(&b, "API Key: %s\n\n", a.keyInput.View())

	connectButtonStyle := styles.ButtonStyle
	if a.connected.Load() {
		connectButtonStyle = styles.ButtonDisabledStyle
	}
	if a.focusedIndex == connectButtFocusIndex {
		connectButtonStyle = styles.SelectedButton(connectButtonStyle)
	}
	fmt.Fprintf(&b, "%s\n\n", connectButtonStyle.Render("Connect"))

	prevStyle := styles.ButtonStyle
	if a.focusedIndex == prevButtonFocusIndex {
		prevStyle = styles.SelectedButton(prevStyle)
	}
	fmt.Fprintf(&b, "%s", prevStyle.Render("Prev"))

	b.WriteRune(' ')

	nextStyle := styles.ButtonStyle
	if !a.CanProceed() {
		nextStyle = styles.ButtonDisabledStyle
	}
	if a.focusedIndex == nextButtonFocusIndex {
		nextStyle = styles.SelectedButton(nextStyle)
	}
	fmt.Fprintf(&b, "%s", nextStyle.Render("Next"))

	b.WriteString("\n\n")

	if a.errorMsg != "" {
		b.WriteString(styles.ErrText.Render("❌ Error: "+a.errorMsg) + "\n")
	}

	if a.connected.Load() {
		b.WriteString("✓ Connected successfully!\n")
	}

	return tea.NewView(b.String())
}

func (a *ApiKeyAuthPage) CanProceed() bool {
	return a.connected.Load()
}

func (a *ApiKeyAuthPage) handleNavigation(key string) {
	switch key {
	case "tab":
		a.focusedIndex = (a.focusedIndex + 1) % 5
	case "shift+tab":
		a.focusedIndex = (a.focusedIndex - 1 + 5) % 5
	case "up":
		a.focusedIndex = (a.focusedIndex - 1 + 5) % 5
	case "down":
		a.focusedIndex = (a.focusedIndex + 1) % 5
	}

	a.addrInput.Blur()
	a.keyInput.Blur()

	switch a.focusedIndex {
	case addrInputFocusIndex:
		a.addrInput.Focus()
	case keyInputFocusIndex:
		a.keyInput.Focus()
	}
}

func (a *ApiKeyAuthPage) connect() tea.Cmd {
	if a.connected.Load() {
		return nil
	}

	client, err := rest.NewClient(a.addrInput.Value())
	if err != nil {
		return nil
	}

	client.SetApiKey(a.keyInput.Value())

	return tea.Cmd(func() tea.Msg {
		if _, err := client.ListSpaces(context.Background()); err != nil {
			return authResultMsg{success: false, err: fmt.Errorf("invalid API key: %w", err)}
		}

		a.setClientOnce.Do(func() {
			a.connected.Store(true)
			a.appAuthState.SetAPIAddress(a.addrInput.Value())
			a.appAuthState.SetAPIKey(a.keyInput.Value())
		})

		return authResultMsg{success: true, err: nil}
	})
}
