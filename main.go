package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
	
help: Displays a help message
exit: Exit the Pokedex`)
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display's help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "closes the program",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		val, ok := commands[scanner.Text()]
		if !ok {
			fmt.Println("Invaild input. \ncall help for more info")
			continue
		}
		val.callback()
	}
}
