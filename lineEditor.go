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

import (
	"github.com/peterh/liner"
)

// LineEditor wraps the basic Prompt method
type LineEditor interface {
	Prompt(string) (string, error)
}

// DefaultLineEditor is a concrete implementation of LineEditor that uses
// github.com/peterh/liner as the line editor
type DefaultLineEditor struct {
	liner *liner.State
}

// Prompt will prompt the user with the prompt string, collect the response and
// return it.  If the upstream liner.Prompt function succeeds, then the
// response is added to the history.  The collected string and any associated
// error is returned
func (d *DefaultLineEditor) Prompt(prompt string) (string, error) {
	str, err := d.liner.Prompt(prompt)
	if err == nil {
		d.liner.AppendHistory(str)
	}
	return str, err
}

// Close returns the terminal to the original state.  This includes taking the
// terminal out of raw mode and turning echo back on
func (d *DefaultLineEditor) Close() error {
	return d.liner.Close()
}

// NewDefaultLineEditor returns a fully initialized line editor that includes
// autocompletion and history
func NewDefaultLineEditor(commands CommandMap) *DefaultLineEditor {
	l := liner.NewLiner()
	l.SetTabCompletionStyle(liner.TabPrints)
	completer := newCompleter(commands)
	l.SetWordCompleter(completer.complete)
	return &DefaultLineEditor{
		liner: l,
	}
}
