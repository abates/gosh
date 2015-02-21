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
	"sort"
	"strings"
	"unicode"
)

type Completer struct {
	topLevelCommands CommandMap
}

func NewCompleter(commands CommandMap) *Completer {
	return &Completer{commands}
}

func (completer Completer) Complete(line string) (c []string) {
	prefix := ""
	fields := strings.Fields(line)
	/* We need to make sure that there are empty fields
	 * in the event of a blank line, or a line that ends
	 * in a space.  Otherwise, there is nothing to attempt
	 * to match on below
	 */
	if len(fields) == 0 {
		fields = []string{""}
	} else if unicode.IsSpace(rune(line[len(line)-1])) {
		fields = append(fields, "")
	}

	commands := completer.topLevelCommands
	for _, field := range fields {
		completions := commands.getCompletions(field)
		for completion, command := range completions {
			/* If it is an exact match then
			 * continue to the next field
			 */
			if field == completion {
				prefix = prefix + completion + " "
				if treeCommand, ok := command.(TreeCommand); ok {
					commands = treeCommand.SubCommands()
					break
				} else if completable, ok := command.(Completable); ok {
					nextCompletions := completable.Completions()
					commands = make(CommandMap, len(nextCompletions))
					for _, nextCompletion := range nextCompletions {
						commands[nextCompletion] = command
					}
				}
			} else {
				c = append(c, prefix+completion)
			}
		}
	}
	sort.Strings(c)
	return
}
