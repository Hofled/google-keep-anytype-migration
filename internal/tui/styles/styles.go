package styles

import "charm.land/lipgloss/v2"

var (
	DisabledText = lipgloss.NewStyle().Strikethrough(true).Faint(true)
	ErrText      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	buttonBackground         = lipgloss.Color("#ff58c4")
	buttonDisabledBackground = lipgloss.Color("#b2b2b2")

	baseButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Padding(0, 3)

	ButtonStyle         = baseButtonStyle.Background(buttonBackground)
	ButtonDisabledStyle = baseButtonStyle.Background(buttonDisabledBackground).Strikethrough(true)
)

func SelectedButton(buttonStyle lipgloss.Style) lipgloss.Style {
	newBackground := lipgloss.Lighten(buttonStyle.GetBackground(), 0.1)
	return buttonStyle.Background(newBackground).Underline(true)
}
