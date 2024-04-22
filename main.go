package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeApi "github.com/djmarkymark007/gopdex/internal/PokeApi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*configCommand) error
}

type configCommand struct {
	Next  *string
	Prevs *string
}

func commandHelp(_ *configCommand) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
	
help: Displays a help message
exit: Exit the Pokedex`)
	return nil
}

func commandExit(_ *configCommand) error {
	os.Exit(0)
	return nil
}

func commandMap(config *configCommand) error {
	if config.Next != nil {
		fmt.Println("No more maps")
		return nil
	}
	data, err := pokeApi.GetLocation(*config.Next)
	if err != nil {
		fmt.Println("failed to get maps")
		return err
	}
	config.Next = data.Next
	config.Prevs = data.Previous
	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(config *configCommand) error {
	if config.Prevs == nil {
		fmt.Println("No prevs maps")
		return nil
	}
	data, err := pokeApi.GetLocation(*config.Prevs)
	if err != nil {
		fmt.Println("failed to get maps")
		return err
	}
	config.Next = data.Next
	config.Prevs = data.Previous
	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
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
		"map": {
			name:        "map",
			description: "get the next 20 location areas in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "get the prevs 20 location areas in the pokemon world",
			callback:    commandMapB,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	baseLocationUrl := "https://pokeapi.co/api/v2/location/"
	config := configCommand{
		Next:  &baseLocationUrl,
		Prevs: nil,
	}

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		val, ok := commands[strings.ToLower(scanner.Text())]
		if !ok {
			fmt.Println("Invaild input. \ncall help for more info")
			continue
		}
		val.callback(&config)
	}
}
