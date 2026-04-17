package challenge

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
)

type CodeModel struct {
	*models.ModelInitOnce

	codeInput textinput.Model

	focusIndex int

	challengeId string
	address     string
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
			case 1:
				return cm, cm.createApiKey()
			}
		}
	case ChallengeIdMsg:
		cm.challengeId = msg.ChallengeId
		cm.address = msg.Address
	}

	switch cm.focusIndex {
	case 0:
		cm.codeInput, cmd = cm.codeInput.Update(msg)
	}

	return cm, cmd
}

func (cm *CodeModel) View() tea.View {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Challenge Code: %s\n", cm.codeInput.View()))

	b.WriteString("\n")

	connectLabel := "Connect"
	if cm.focusIndex == 1 {
		connectLabel = fmt.Sprintf("[%s]", connectLabel)
	}
	b.WriteString(fmt.Sprintf("%s\n", connectLabel))

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
	case 0:
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
			// TODO
			return nil
		}

		return ApiKeyMsg{
			ApiKey: apiKeyRes.ApiKey,
		}
	}
}
