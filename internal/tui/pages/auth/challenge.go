package auth

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth/challenge"
)

type viewState uint

const (
	initView viewState = iota
	codeView
)

type ChallengeAuthPage struct {
	*models.ModelInitOnce
	*models.PageIds

	appPageState state.AppPageStater
	appAuthState state.AppAuthStater

	currentView viewState

	initChallenge *challenge.InitModel
	challengeCode *challenge.CodeModel

	focusedIndex int
}

func NewChallengeAuthPage(appAuthStater state.AppAuthStater, appPageState state.AppPageStater) (*ChallengeAuthPage, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	p := &ChallengeAuthPage{
		PageIds:       pageIds,
		appAuthState:  appAuthStater,
		appPageState:  appPageState,
		currentView:   initView,
		initChallenge: challenge.NewInitModel(),
		challengeCode: challenge.NewCodeModel(),
		focusedIndex:  0,
	}

	p.ModelInitOnce = models.NewModelInitOnce(p)

	return p, nil
}

func (cap *ChallengeAuthPage) Init() tea.Cmd {
	return tea.Batch(cap.initChallenge.Init(), cap.challengeCode.Init())
}

func (cap *ChallengeAuthPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case challenge.ChallengeIdMsg:
		cap.currentView = codeView
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*challenge.CodeModel)
		break
	}

	switch cap.currentView {
	case initView:
		var m tea.Model
		m, cmd = cap.initChallenge.Update(msg)
		cap.initChallenge = m.(*challenge.InitModel)
		break
	case codeView:
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*challenge.CodeModel)
		break
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
}

func (cap *ChallengeAuthPage) View() tea.View {
	var b strings.Builder

	b.WriteString("Challenge Authentication\n\n")

	var subView string

	switch cap.currentView {
	case initView:
		subView = cap.initChallenge.View().Content
		break
	case codeView:
		subView = cap.challengeCode.View().Content
		break
	}

	b.WriteString(subView)

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
