package models

import (
	"sync"

	tea "charm.land/bubbletea/v2"
)

type ModelInitOncer interface {
	InitOnce() tea.Cmd
}

type InitOnceFunc func() tea.Cmd

type ModelInitOnce struct {
	initOnceFunc InitOnceFunc
}

func NewModelInitOnce(model tea.Model) *ModelInitOnce {
	return &ModelInitOnce{
		initOnceFunc: createInitOnceFunc(model),
	}
}

func (mic *ModelInitOnce) InitOnce() tea.Cmd {
	return mic.initOnceFunc()
}

func createInitOnceFunc(model tea.Model) InitOnceFunc {
	return sync.OnceValue(func() tea.Cmd {
		return model.Init()
	})
}
