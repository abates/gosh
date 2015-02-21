package gosh

import (
	"errors"
)

var ErrDuplicateCommand = errors.New("Command already exists")
var ErrNilCallback = errors.New("Cannot assign nil callback functions")
var ErrNilLineEditor = errors.New("Cannot assign a nil line editor")
var ErrNilPrompt = errors.New("Cannot assign a nil prompt")
var ErrNilWriter = errors.New("Cannot assign a nil writer")
var ErrNoMatchingCommand = errors.New("No matching command")
