package models

import (
	tea "charm.land/bubbletea/v2"
)

type ModelOnceIniter interface {
	InitOnce() tea.Cmd
}
