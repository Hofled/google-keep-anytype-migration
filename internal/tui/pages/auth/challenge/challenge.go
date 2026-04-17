package challenge

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/styles"
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

type keyMap struct {
	Back key.Binding
	Tab  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Tab}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Tab},
	}
}

type ChallengeAuthPage struct {
	*models.ModelInitOnce
	*models.PageIds

	appPageState state.AppPageStater
	appAuthState state.AppAuthStater

	currentSubView viewState

	initChallenge *InitModel
	challengeCode *CodeModel

	subViewFocused bool

	focusIndex int

	connected bool

	help       help.Model
	keyMapping keyMap
}

func NewChallengeAuthPage(appAuthStater state.AppAuthStater, appPageState state.AppPageStater) (*ChallengeAuthPage, error) {
	keyMapping := keyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc ▣", "return to challenge view"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("⇆ / tab", "change focus between view and nav"),
		),
	}

	pageIds, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	p := &ChallengeAuthPage{
		PageIds:        pageIds,
		appAuthState:   appAuthStater,
		appPageState:   appPageState,
		currentSubView: initView,
		initChallenge:  NewInitModel(),
		challengeCode:  NewCodeModel(),
		focusIndex:     0,
		subViewFocused: true,
		help:           help.New(),
		keyMapping:     keyMapping,
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
		pressedKey := msg.String()
		if key.Matches(msg.Key(), cap.keyMapping.Back) {
			cap.currentSubView = initView
		}

		cap.handleViewFocus(pressedKey)

		if !cap.subViewFocused {
			cap.handleNavigation(pressedKey)
			if pressedKey == "enter" {
				switch cap.focusIndex {
				case nextFocusIndex:
					if cap.connected {
						cap.appPageState.NextPage()
					}
					return cap, nil
				case prevFocusIndex:
					cap.appPageState.PrevPage()
					return cap, nil
				}
			}

			return cap, nil
		}
	case ChallengeIdMsg:
		cap.currentSubView = codeView
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*CodeModel)
	case ApiKeyMsg:
		cap.appAuthState.SetAPIKey(msg.ApiKey)
		cap.connected = true
	}

	switch cap.currentSubView {
	case initView:
		var m tea.Model
		m, cmd = cap.initChallenge.Update(msg)
		cap.initChallenge = m.(*InitModel)
	case codeView:
		var m tea.Model
		m, cmd = cap.challengeCode.Update(msg)
		cap.challengeCode = m.(*CodeModel)
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
		cap.focusIndex = (cap.focusIndex + 1) % 2
	case "left":
		cap.focusIndex = (cap.focusIndex - 1 + 2) % 2
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
	if !cap.subViewFocused && cap.focusIndex == prevFocusIndex {
		prevLabel = fmt.Sprintf("[%s]", prevLabel)
	}
	b.WriteString(fmt.Sprintf("%s", prevLabel))

	b.WriteRune(' ')

	nextLabel := "Next"
	if !cap.connected {
		nextLabel = styles.DisabledText.Render(nextLabel)
	}
	if !cap.subViewFocused && cap.focusIndex == nextFocusIndex {
		nextLabel = fmt.Sprintf("[%s]", nextLabel)
	}
	b.WriteString(fmt.Sprintf("%s", nextLabel))

	b.WriteString("\n")

	bString := b.String()

	helpView := cap.help.View(cap.keyMapping)
	height := 8 - strings.Count(bString, "\n") - strings.Count(helpView, "\n")

	return tea.NewView(bString + strings.Repeat("\n", height) + helpView)
}
