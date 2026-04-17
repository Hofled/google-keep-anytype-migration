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

const (
	nextFocusIndex int = iota
	prevFocusIndex
)

type ChallengeAuthPage struct {
	*models.ModelInitOnce
	*models.PageIds

	appPageState state.AppPageStater
	appAuthState state.AppAuthStater

	currentSubView viewState

	initChallenge *challenge.InitModel
	challengeCode *challenge.CodeModel

	subViewFocused bool

	focusedIndex int
}

func NewChallengeAuthPage(appAuthStater state.AppAuthStater, appPageState state.AppPageStater) (*ChallengeAuthPage, error) {
	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	p := &ChallengeAuthPage{
		PageIds:        pageIds,
		appAuthState:   appAuthStater,
		appPageState:   appPageState,
		currentSubView: initView,
		initChallenge:  challenge.NewInitModel(),
		challengeCode:  challenge.NewCodeModel(),
		focusedIndex:   0,
		subViewFocused: true,
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
	case tea.KeyPressMsg:
		key := msg.String()
		cap.handleViewFocus(key)

		if !cap.subViewFocused {
			cap.handleNavigation(key)
			if key == "enter" {
				switch cap.focusedIndex {
				case nextFocusIndex:
					cap.appPageState.NextPage()
					return cap, nil
				case prevFocusIndex:
					cap.appPageState.PrevPage()
					return cap, nil
				}
			}

			return cap, nil
		}
	case challenge.ChallengeIdMsg:
		cap.currentSubView = codeView
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*challenge.CodeModel)
	}

	switch cap.currentSubView {
	case initView:
		var m tea.Model
		m, cmd = cap.initChallenge.Update(msg)
		cap.initChallenge = m.(*challenge.InitModel)
	case codeView:
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*challenge.CodeModel)
	}

	return cap, cmd
}

func (cap *ChallengeAuthPage) handleViewFocus(key string) {
	switch key {
	case "tab", "shift+tab":
		cap.subViewFocused = !cap.subViewFocused
	}
}

func (cap *ChallengeAuthPage) handleNavigation(key string) {
	switch key {
	case "right":
		cap.focusedIndex = (cap.focusedIndex + 1) % 2
	case "left":
		cap.focusedIndex = (cap.focusedIndex - 1 + 2) % 2
	}
}

func (cap *ChallengeAuthPage) View() tea.View {
	var b strings.Builder

	b.WriteString("Challenge Authentication\n\n")

	var subView string

	switch cap.currentSubView {
	case initView:
		subView = cap.initChallenge.View().Content
	case codeView:
		subView = cap.challengeCode.View().Content
	}

	b.WriteString(subView)

	prevLabel := "Prev"
	if !cap.subViewFocused && cap.focusedIndex == prevFocusIndex {
		prevLabel = fmt.Sprintf("[%s]", prevLabel)
	}
	b.WriteString(fmt.Sprintf("%s", prevLabel))

	b.WriteRune(' ')

	nextLabel := "Next"
	if !cap.subViewFocused && cap.focusedIndex == nextFocusIndex {
		nextLabel = fmt.Sprintf("[%s]", nextLabel)
	}
	b.WriteString(fmt.Sprintf("%s", nextLabel))

	b.WriteString("\n")

	return tea.NewView(b.String())
}
