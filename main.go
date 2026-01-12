package main

import (
	"fmt"
	"os"
	"wnetctl/command"
)

func main() {
	argv := os.Args[1:]
	cmd := handleCommandLine(argv)
	if cmd.HelpRequested() {
		fmt.Println(cmd.HelpMessage())
	} else {
		if err := cmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
	}
}

func handleCommandLine(argv []string) command.Command {
	switch argv[0] {
	case "site":
		return command.GetSiteCommand(argv[1:])
	case "ap":
		return command.GetApCommand(argv[1:])
	case "ssid":
		return command.GetSsidCommand(argv[1:])
	case "device":
		return command.GetDeviceCommand(argv[1:])
	case "station":
		return command.Help(true)
	case "help":
		return command.Help(true)
	}
	return command.Help(true)
}
