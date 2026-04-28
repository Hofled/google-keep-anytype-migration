package styles

import "charm.land/lipgloss/v2"

var (
	DisabledText = lipgloss.NewStyle().Strikethrough(true).Faint(true)
	ErrText      = lipgloss.NewStyle().Foreground(lipgloss.Red)
)
