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
	completer := &Completer{commands}
	return completer
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
				commands = command.SubCommands()
				break
			} else {
				c = append(c, prefix+completion)
			}
		}
	}
	sort.Strings(c)
	return
}
