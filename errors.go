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
	"errors"
)

var (
	// ErrDefaultPrompter indicates that the Prompt is not a DefaultPrompt so the
	// prompter function cannot be overridden
	ErrDefaultPrompter = errors.New("can only set the prompter on the DefaultPrompt")

	// ErrDuplicateCommand indicates a command with the same name already exists
	// in the CommandMap
	ErrDuplicateCommand = errors.New("command already exists")

	// ErrNilCallback indicates that a callback function was set to nil
	ErrNilCallback = errors.New("cannot assign nil callback functions")

	// ErrNilLineEditor indicates that the Prompt's line editor was set to nil
	ErrNilLineEditor = errors.New("cannot assign a nil line editor")

	// ErrNilPrompt indicates that the Shell's Prompt was set to nil
	ErrNilPrompt = errors.New("cannot assign a nil prompt")

	// ErrNilPrompter indicates that the Prompt's prompter was set to nil
	ErrNilPrompter = errors.New("cannot assign a nil prompter")

	// ErrNilWriter indicates the shell's writer was set to nil
	ErrNilWriter = errors.New("cannot assign a nil writer")

	// ErrNoMatchingCommand indicates a matching command could not be found in the CommandMap
	ErrNoMatchingCommand = errors.New("no matching command")
)
