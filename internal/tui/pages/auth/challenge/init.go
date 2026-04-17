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

type InitModel struct {
	*models.ModelInitOnce

	addrInput textinput.Model

	createChallengeErr error

	focusIndex int
}

type ChallengeIdMsg struct {
	ChallengeId string
	Address     string
}

func NewInitModel() *InitModel {
	addrInput := textinput.New()
	addrInput.SetValue("http://localhost:31009")
	addrInput.Placeholder = "http://localhost:31009"
	addrInput.Focus()
	addrInput.SetWidth(50)

	return &InitModel{
		addrInput: addrInput,
	}
}

func (im *InitModel) Init() tea.Cmd {
	return textinput.Blink
}

func (im *InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			im.handleNavigation(msg.String())
			return im, nil

		case "enter":
			switch im.focusIndex {
			case 1:
				return im, im.startChallenge()
			}
		}
	}

	switch im.focusIndex {
	case 0:
		im.addrInput, cmd = im.addrInput.Update(msg)
	}

	return im, cmd
}

func (im *InitModel) View() tea.View {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("API Address: %s\n", im.addrInput.View()))

	b.WriteString("\n")

	challengeLabel := "Challenge"
	if im.focusIndex == 1 {
		challengeLabel = fmt.Sprintf("[%s]", challengeLabel)
	}
	b.WriteString(fmt.Sprintf("%s\n", challengeLabel))

	if im.createChallengeErr != nil {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Error starting challenge: %s\n", im.createChallengeErr))
	}

	return tea.NewView(b.String())
}

func (im *InitModel) handleNavigation(key string) {
	switch key {
	case "tab", "down":
		im.focusIndex = (im.focusIndex + 1) % 2
	case "shift+tab", "up":
		im.focusIndex = (im.focusIndex - 1 + 2) % 2
	}

	im.addrInput.Blur()

	switch im.focusIndex {
	case 0:
		im.addrInput.Focus()
	}
}

func (im *InitModel) startChallenge() tea.Cmd {
	im.createChallengeErr = nil

	restClient, err := rest.NewClient(im.addrInput.Value())
	if err != nil {
		im.createChallengeErr = err
		return nil
	}

	return func() tea.Msg {
		challengeRes, err := restClient.CreateChallenge(context.Background())
		if err != nil {
			im.createChallengeErr = err
			return nil
		}

		return ChallengeIdMsg{
			ChallengeId: challengeRes.ChallengeId,
			Address:     im.addrInput.Value(),
		}
	}
}
