package gosh

import ()

type Closeable interface {
	Close() error
}

type Prompt interface {
	NextResponse() (string, error)
}

type DefaultPrompt struct {
	prompter   func() string
	lineEditor LineEditor
	commands   CommandMap
}

func NewDefaultPrompt(commands CommandMap) *DefaultPrompt {
	p := DefaultPrompt{
		func() string {
			return "> "
		},
		NewDefaultLineEditor(),
		commands,
	}
	return &p
}

func (p *DefaultPrompt) SetPrompter(prompter func() string) error {
	if prompter == nil {
		return ErrNilCallback
	}
	p.prompter = prompter
	return nil
}

func (p *DefaultPrompt) SetLineEditor(lineEditor LineEditor) error {
	if lineEditor == nil {
		return ErrNilLineEditor
	}
	p.lineEditor = lineEditor
	return nil
}

func (p *DefaultPrompt) NextResponse() (string, error) {
	return p.lineEditor.Prompt(p.prompter())
}

func (p *DefaultPrompt) Close() error {
	if lineEditor, ok := p.lineEditor.(Closeable); ok {
		return lineEditor.Close()
	}
	return nil
}
