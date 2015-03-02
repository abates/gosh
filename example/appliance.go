package main

import (
	"fmt"
	"github.com/abates/gosh"
	"net"
	"os"
	"time"
)

type TimeCommand struct{}

func (TimeCommand) Exec() error {
	t := time.Now()
	fmt.Println(t.Format(time.RFC822))
	return nil
}

func interfaceNames() ([]string, error) {
	interfaces, err := net.Interfaces()
	names := make([]string, len(interfaces))
	if err != nil {
		return nil, err
	} else {
		for i, netInterface := range interfaces {
			names[i] = netInterface.Name
		}
	}
	return names, nil
}

type InterfaceCommand struct{}

func (i InterfaceCommand) Completions(substring string) []string {
	names, err := interfaceNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve system interfaces: %v\n", err)
	}
	return names
}

func (i InterfaceCommand) Exec() error {
	for _, name := range os.Args[1:] {
		netInterface, err := net.InterfaceByName(name)
		if err != nil {
			return err
		} else {
			fmt.Printf("Name: %v\n", netInterface.Name)
			addresses, err := netInterface.Addrs()
			if err != nil {
				return err
			}
			for _, address := range addresses {
				fmt.Printf("      %v\n", address)
			}
		}
	}
	return nil
}

type InterfacesCommand struct{}

func (i InterfacesCommand) Exec() error {
	names, err := interfaceNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Fprintf(os.Stdout, "%v\n", name)
	}
	return nil
}

var commands = gosh.CommandMap{
	"show": gosh.NewTreeCommand(gosh.CommandMap{
		"interface":  InterfaceCommand{},
		"interfaces": InterfacesCommand{},
		"time":       TimeCommand{},
	}),
}

func main() {
	shell := gosh.NewShell(commands)
	shell.Exec()
}
