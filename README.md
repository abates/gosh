# gosh version 0.0.1

Gosh is a simple package designed to make creating command line shells in go a
little bit easier.  The concept is simple: have a command prompt with
auto-completion and history that can execute commands written in go.  The
commands can be hierarchical, similar to the CLI in common network operating
systems.  The commands can also be more like traditional OS commands, with no
hierarchy.

## Examples

A simple example of a single command shell:

```go
package main

import (
  "fmt"
  "github.com/abates/gosh"
)

type cmd string

func (c cmd) Exec() error {
  fmt.Printf("Executing %s\n", string(c))
  return nil
}

var commands = gosh.CommandMap{
  "cmd": cmd("My Command!"),
}

func main() {
  shell := gosh.NewShell(commands)
  shell.Exec()
}
```

Any arguments that follow the command on the prompt are passed into the Exec
method by way of the os.Args string slice.  The shell will initialize with a
default line editor that implements history, auto-completion and the prompt
string.  

The default shell can be supplied with a customized prompter:

```go
func main() {
  shell := gosh.NewShell(commands)
  shell.SetPrompter(func() string {
    return "Custom Prompt> "
  })
  shell.Exec()
}
```

Commands can optionally specify auto-completion candidates for their arguments:
```go
type cmd string

func (c cmd) Exec() error {
  fmt.Printf("Executing %s\n", string(c))
  return nil
}

func (c cmd) Completions(field string) []string {
  return []string{"arg1", "arg2", "arg3"}
}
```

Command hierarchies can be created with the CommandTree struct:
```go
var commands = gosh.CommandMap{
  "show": gosh.NewTreeCommand(gosh.CommandMap{
    "interface":  InterfaceCommand{},
    "interfaces": InterfacesCommand{},
    "time":       TimeCommand{},
  }),
}
```

## Documentation
https://godoc.org/github.com/abates/gosh

