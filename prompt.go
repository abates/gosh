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

// Closeable objects implement a Close method.
//
// If the prompt detects that the line editor implements Closeable then it will
// call Close upon completing the prompt cycle
type Closeable interface {
	Close() error
}

// Prompt receivers implement a NextResponse method
//
// NextResponse should return any response provided by the user as well as any
// error encountered during the prompting
type Prompt interface {
	NextResponse() (string, error)
}

// Prompter returns the prompt string used in NextResponse
//
// Prompter is a function that is called to generate the prompt string that
// will preced user input.  For instance, if a Prompter returns the string "> "
// then all input lines will begin with that string and user input starts just
// after the space following the ">"
type Prompter func() string

// DefaultPrompt is a concrete implementation of Prompt
//
// DefaultPrompt includes a basi "> " prompt and the DefaultLineEditor which
// provides tab completion and command history
type DefaultPrompt struct {
	prompter   Prompter
	lineEditor LineEditor
}

// NewDefaultPrompt returns a fully initialized DefaultPrompt.
//
// When NextResponse is called on the DefaultPrompt the prompt will be "> ".
// The DefaultLineEditor is used so tab completions and history is available
func NewDefaultPrompt(commands CommandMap) *DefaultPrompt {
	p := DefaultPrompt{
		func() string {
			return "> "
		},
		NewDefaultLineEditor(commands),
	}
	return &p
}

// SetPrompter will allow overriding the default prompter.
//
// If the prompter argument is nil then the prompter is not overridden and the
// ErrNilCallback is returned.  Otherwise the prompter is set to whatever
// function is provided
func (p *DefaultPrompt) SetPrompter(prompter func() string) error {
	if prompter == nil {
		return ErrNilCallback
	}
	p.prompter = prompter
	return nil
}

// SetLineEditor allows overriding the default line editor
//
// The DefaultPrompt uses github.com/peterh/liner as the default line editor.
// This can be overridden with a different LineEditor implementation.  If
// SetLineEditor is called with a nil LineEditor then the ErrNilLineEditor
// error is returned
func (p *DefaultPrompt) SetLineEditor(lineEditor LineEditor) error {
	if lineEditor == nil {
		return ErrNilLineEditor
	}
	p.lineEditor = lineEditor
	return nil
}

// NextResponse prompts the user and returns the uer's input
func (p *DefaultPrompt) NextResponse() (string, error) {
	return p.lineEditor.Prompt(p.prompter())
}

// Close closes the line editor (if it is closeable)
func (p *DefaultPrompt) Close() error {
	if lineEditor, ok := p.lineEditor.(Closeable); ok {
		return lineEditor.Close()
	}
	return nil
}
