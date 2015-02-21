/**
 * Copyright 2015 Andrew Bates
 *
 * Licensed under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package gosh

import ()

type Closeable interface {
	Close() error
}

type Prompt interface {
	NextResponse() (string, error)
}

type Prompter func() string

type DefaultPrompt struct {
	prompter   Prompter
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
