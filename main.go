package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	pokeApi "github.com/djmarkymark007/gopdex/internal/pokeApi"
	"github.com/djmarkymark007/gopdex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*configCommand, []string) error
}

type configCommand struct {
	LocationUrl string
	Next        *string
	Prevs       *string
}

func commandHelp(_ *configCommand, _ []string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
	
help: Displays a help message
exit: Exit the Pokedex`)
	return nil
}

func commandExit(_ *configCommand, _ []string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *configCommand, _ []string) error {
	start := time.Now()
	if config.Next == nil {
		fmt.Println("No more maps")
		return nil
	}
	data, err := getLocation(*config.Next)
	if err != nil {
		fmt.Println("failed to get maps")
		return err
	}
	config.Next = data.Next
	config.Prevs = data.Previous
	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
	end := time.Now()
	fmt.Printf("took: %v\n", end.Sub(start))
	return nil
}

func commandMapB(config *configCommand, _ []string) error {
	start := time.Now()
	if config.Prevs == nil {
		fmt.Println("No prevs maps")
		return nil
	}
	data, err := getLocation(*config.Prevs)
	if err != nil {
		fmt.Println("failed to get maps")
		return err
	}
	config.Next = data.Next
	config.Prevs = data.Previous
	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
	end := time.Now()
	fmt.Printf("took: %v\n", end.Sub(start))
	return nil
}

func commandExpore(config *configCommand, area_name []string) error {
	if len(area_name) != 1 {
		fmt.Print("Explore has only one arg (location)\n")
		return errors.New("too many args")
	}
	url := config.LocationUrl + area_name[0]
	data, err := getLocationPokemon(url)
	if err != nil {
		return err
	}

	for _, encounter := range data.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func getLocationPokemon(url string) (pokeApi.LocationPokemon, error) {
	fmt.Println(url)
	data, err := cache.Get(url)
	var errIn error
	if !err {
		data, errIn = pokeApi.GetLocation(url)
		if errIn != nil {
			fmt.Println("failed to get maps")
			return pokeApi.LocationPokemon{}, errIn
		}
	}

	var pokemon pokeApi.LocationPokemon
	pokemon, errIn = pokeApi.JsonToLocationPokemon(data)
	if errIn != nil {
		fmt.Printf("failed to convert json to location pokemon struct %v\n", errIn.Error())
		return pokeApi.LocationPokemon{}, errIn
	}

	return pokemon, nil
}

func getLocation(url string) (pokeApi.Location, error) {
	fmt.Println(url)
	data, err := cache.Get(url)
	var errIn error
	if !err {
		data, errIn = pokeApi.GetLocation(url)
		if errIn != nil {
			fmt.Println("failed to get maps")
			return pokeApi.Location{}, errIn
		}
	}

	var loc pokeApi.Location
	loc, errIn = pokeApi.JsonToLocation(data)
	if errIn != nil {
		fmt.Println("failed to convert json to location struct")
		return pokeApi.Location{}, errIn
	}
	return loc, nil
}

var cache pokecache.Cache

func main() {
	cache = *pokecache.NewCache(5 * time.Minute)
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
		"explore": {
			name:        "explore",
			description: "get the pokemon for the location",
			callback:    commandExpore,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	baseLocationUrl := "https://pokeapi.co/api/v2/location/"
	LocationPokemonUrl := "https://pokeapi.co/api/v2/location-area/"
	config := configCommand{
		LocationUrl: LocationPokemonUrl,
		Next:        &baseLocationUrl,
		Prevs:       nil,
	}

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := strings.Split(scanner.Text(), " ")
		command := input[0]
		args := input[1:]
		val, ok := commands[strings.ToLower(command)]
		if !ok {
			fmt.Println("Invaild input. \ncall help for more info")
			continue
		}
		val.callback(&config, args)
	}
}
