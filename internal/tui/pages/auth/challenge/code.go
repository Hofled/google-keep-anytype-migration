package challenge

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/styles"
	"github.com/Hofled/go-google-keep-anytype-migration/pkg/tui/widgets"
)

const (
	codeInputFocusIndex = iota
	connectButtFocusIndex
)

type CodeModel struct {
	*models.ModelInitOnce
	widgets.FocusableWidget

	codeInput textinput.Model

	focusIndex int

	challengeId string
	address     string

	done bool
}

type ApiKeyMsg struct {
	ApiKey string
}

func NewCodeModel() *CodeModel {
	codeInput := textinput.New()
	codeInput.Placeholder = "1234"
	codeInput.SetWidth(4)
	codeInput.CharLimit = 4
	codeInput.Focus()

	return &CodeModel{
		codeInput: codeInput,
	}
}

func (cm *CodeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cm *CodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "down":
			cm.handleNavigation(msg.String())
			return cm, nil

		case "enter":
			switch cm.focusIndex {
			case connectButtFocusIndex:
				return cm, cm.createApiKey()
			}
		}
	case ChallengeIdMsg:
		cm.challengeId = msg.ChallengeId
		cm.address = msg.Address
	}

	switch cm.focusIndex {
	case codeInputFocusIndex:
		cm.codeInput, cmd = cm.codeInput.Update(msg)
	}

	return cm, cmd
}

func (cm *CodeModel) View() tea.View {
	var b strings.Builder

	fmt.Fprintf(&b, "Challenge Code: %s\n\n", cm.codeInput.View())

	connectButtonStyle := styles.ButtonGrayedOutStyle
	if cm.done {
		connectButtonStyle = styles.ButtonDisabledStyle
	} else if cm.Focused() {
		connectButtonStyle = styles.ButtonStyle
		if cm.focusIndex == connectButtFocusIndex {
			connectButtonStyle = styles.SelectedButton(connectButtonStyle)
		}
	}

	fmt.Fprintf(&b, "%s\n", connectButtonStyle.Render("Connect"))

	return tea.NewView(b.String())
}

func (cm *CodeModel) handleNavigation(key string) {
	switch key {
	case "down":
		cm.focusIndex = (cm.focusIndex + 1) % 2
	case "up":
		cm.focusIndex = (cm.focusIndex - 1 + 2) % 2
	}

	cm.codeInput.Blur()

	switch cm.focusIndex {
	case codeInputFocusIndex:
		cm.codeInput.Focus()
	}
}

func (cm *CodeModel) createApiKey() tea.Cmd {
	restClient, err := rest.NewClient(cm.address)
	if err != nil {
		return nil
	}

	return func() tea.Msg {
		apiKeyRes, err := restClient.CreateApiKey(context.Background(), cm.challengeId, cm.codeInput.Value())
		if err != nil {
			// TODO reflect error in component
			return nil
		}

		cm.done = true

		return ApiKeyMsg{
			ApiKey: apiKeyRes.ApiKey,
		}
	}
}
