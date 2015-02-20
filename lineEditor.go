package gosh

import (
	"github.com/peterh/liner"
)

type LineEditor interface {
	Prompt(string) (string, error)
}

type DefaultLineEditor struct {
	liner *liner.State
}

func (editor DefaultLineEditor) Prompt(prompt string) (string, error) {
	return editor.liner.Prompt(prompt)
}

func NewDefaultLineEditor() LineEditor {
	editor := DefaultLineEditor{
		liner: liner.NewLiner(),
	}

	return editor
}

func (editor DefaultLineEditor) Close() error {
	return editor.liner.Close()
}
