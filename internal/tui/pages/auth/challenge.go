package auth

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

type ChallengeAuthPage struct {
	*models.ModelInitOnce
	*models.PageIds

	appPageState state.AppPageStater
	appAuthState state.AppAuthStater

	addrInput    textinput.Model
	focusedIndex int

	challengeId string
}

func NewChallengeAuthPage(appAuthStater state.AppAuthStater, appPageState state.AppPageStater) (*ChallengeAuthPage, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	addrInput := textinput.New()
	addrInput.SetValue("http://localhost:31009")
	addrInput.Placeholder = "http://localhost:31009"
	addrInput.Focus()
	addrInput.SetWidth(50)

	p := &ChallengeAuthPage{
		PageIds:      pageIds,
		appAuthState: appAuthStater,
		appPageState: appPageState,
		addrInput:    addrInput,
		focusedIndex: 0,
	}

	p.ModelInitOnce = models.NewModelInitOnce(p)

	return p, nil
}

func (cap *ChallengeAuthPage) Init() tea.Cmd {
	return nil
}

func (cap *ChallengeAuthPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			cap.handleNavigation(msg.String())
			return cap, nil
		case "enter":
			switch cap.focusedIndex {
			case 1:
				challengeId, err := cap.challenge()
				if err == nil {
					fmt.Println(challengeId)
				}
			case 2:
				cap.appPageState.NextPage()
				return cap, nil
			case 3:
				cap.appPageState.PrevPage()
				return cap, nil
			}
		}
	}

	switch cap.focusedIndex {
	case 0:
		cap.addrInput, cmd = cap.addrInput.Update(msg)
	}

	return cap, cmd
}

func (cap *ChallengeAuthPage) handleNavigation(key string) {
	switch key {
	case "tab":
		cap.focusedIndex = (cap.focusedIndex + 1) % 4
	case "shift+tab":
		cap.focusedIndex = (cap.focusedIndex - 1 + 4) % 4
	case "up":
		cap.focusedIndex = (cap.focusedIndex - 1 + 4) % 4
	case "down":
		cap.focusedIndex = (cap.focusedIndex + 1) % 4
	}

	cap.addrInput.Blur()

	switch cap.focusedIndex {
	case 0:
		cap.addrInput.Focus()
	}
}

func (cap *ChallengeAuthPage) View() tea.View {
	var b strings.Builder

	b.WriteString("Challenge Authentication\n\n")

	b.WriteString(fmt.Sprintf("API Address: %s\n", cap.addrInput.View()))

	b.WriteString("\n")

	challengeLabel := "Challenge"
	if cap.focusedIndex == 1 {
		challengeLabel = fmt.Sprintf("[%s]", challengeLabel)
	}
	b.WriteString(fmt.Sprintf("%s\n", challengeLabel))

	nextLabel := "Next"
	if cap.focusedIndex == 2 {
		nextLabel = fmt.Sprintf("[%s]", nextLabel)
	}
	b.WriteString(fmt.Sprintf("%s\n", nextLabel))

	prevLabel := "Prev"
	if cap.focusedIndex == 3 {
		prevLabel = fmt.Sprintf("[%s]", prevLabel)
	}
	b.WriteString(fmt.Sprintf("%s\n", prevLabel))

	return tea.NewView(b.String())
}

func (cap *ChallengeAuthPage) challenge() (string, error) {
	restClient, err := rest.NewClient(cap.addrInput.Value())
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	challengeRes, err := restClient.CreateChallenge(ctx)

	return challengeRes.ChallengeId, err
}
