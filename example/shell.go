package main

import (
	"errors"
	"fmt"
	"github.com/abates/gosh"
	"os"
	"path"
)

var err error
var cwd string

type cd struct{}

func (this cd) Completions(arg string) []string {
	dir, _ := path.Split(arg)
	lsdir := ""
	if dir == "" {
		lsdir = cwd
	} else {
		if string(dir[0]) == "/" {
			lsdir = dir
		} else {
			lsdir = cwd + dir
		}
	}

	f, err := os.Open(lsdir)
	defer f.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		return []string{}
	}
	names, err := f.Readdirnames(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		return []string{}
	}

	candidates := make([]string, len(names))
	i := 0
	for _, name := range names {
		if isDir(lsdir + name) {
			candidates[i] += dir + name + "/"
		} else {
			candidates[i] += dir + name + " "
		}
		i++
	}
	return candidates
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (this cd) Exec(arguments []string) error {
	nextDir := ""
	if len(arguments) == 0 {
		home := os.Getenv("HOME")
		if home != "" {
			cwd = home
		}
		return nil
	}
	dir := arguments[0]
	if dir[0] == '/' {
		nextDir = string(dir[0])
	} else {
		nextDir = cwd + "/" + dir
	}
	if isDir(nextDir) {
		cwd = path.Clean(nextDir)
		return nil
	}
	return errors.New(fmt.Sprintf("%v is not a directory", dir))
}

type ls struct{}

func (this ls) Exec(arguments []string) error {
	f, err := os.Open(cwd)
	defer f.Close()

	if err != nil {
		return err
	}
	names, err := f.Readdirnames(0)
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

var commands = gosh.CommandMap{
	"cd": cd{},
	"ls": ls{},
}

func main() {
	cwd = "/"
	shell := gosh.NewShell(commands)
	shell.SetPrompter(func() string {
		return fmt.Sprintf("%v> ", cwd)
	})
	shell.Exec()
}
