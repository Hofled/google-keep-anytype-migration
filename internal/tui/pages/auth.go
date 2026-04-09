package pages

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

var disabledTextStyle = lipgloss.NewStyle().Strikethrough(true).Faint(true)

type AuthPage struct {
	*models.ModelInitOnce

	addrInput    textinput.Model
	keyInput     textinput.Model
	errorMsg     string
	connected    bool
	focusedIndex int

	appAuthState state.AppAuthStater
	appViewState state.AppViewStater
}

type authResultMsg struct {
	success bool
	err     error
}

func NewAuthPage(appAuthState state.AppAuthStater, appViewState state.AppViewStater) *AuthPage {
	addrInput := textinput.New()
	addrInput.SetValue("https://localhost:31009")
	addrInput.Placeholder = "https://localhost:31009"
	addrInput.Focus()
	addrInput.SetWidth(50)

	keyInput := textinput.New()
	keyInput.Placeholder = "Your API Key"
	keyInput.SetWidth(50)

	authPage := &AuthPage{
		addrInput:    addrInput,
		keyInput:     keyInput,
		focusedIndex: 0,
		appAuthState: appAuthState,
		appViewState: appViewState,
	}

	authPage.ModelInitOnce = models.NewModelInitOnce(authPage)

	return authPage
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
				// Proceed to next page (handled by app)
				return a, nil
			}
		}
	case authResultMsg:
		if msg.success {
			a.connected = true
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
	if a.focusedIndex == 2 {
		connectLabel = "[" + connectLabel + "]"
	}
	b.WriteString(fmt.Sprintf("%s\n", connectLabel))

	nextLabel := "Next"
	if !a.CanProceed() {
		nextLabel = disabledTextStyle.Render("(Next)")
	} else if a.focusedIndex == 3 {
		nextLabel = "[" + nextLabel + "]"
	}
	b.WriteString(fmt.Sprintf("%s\n\n", nextLabel))

	if a.errorMsg != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("❌ Error: "+a.errorMsg) + "\n")
	}

	if a.connected {
		b.WriteString("✓ Connected successfully!\n")
	}

	return tea.NewView(b.String())
}

func (a *AuthPage) CanProceed() bool {
	return a.connected
}

func (a *AuthPage) GetData() any {
	return map[string]string{
		"addr": a.addrInput.Value(),
		"key":  a.keyInput.Value(),
	}
}

func (a *AuthPage) NextPage() {
	a.appViewState.NextView()
}

func (a *AuthPage) PrevPage() {
	a.appViewState.PrevView()
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
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		authenticatedClient, err := anytype.AuthWithChallenge(ctx, a.addrInput.Value(), a.keyInput.Value())
		if err != nil {
			return authResultMsg{success: false, err: err}
		}

		a.appAuthState.SetClient(authenticatedClient)

		return authResultMsg{success: true, err: nil}
	})
}
