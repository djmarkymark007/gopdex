package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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
	Pokedex     map[string]pokeApi.Pokemon
	PokemonUrl  string
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
	//TODO(Mark): check the case sinsativity of area_name
	url := config.LocationUrl + area_name[0]

	fmt.Printf("url: %v\n", url)
	data, err := cache.Get(url)
	var errIn error
	if !err {
		data, errIn = pokeApi.CallApiByUrl(url)
		if errIn != nil {
			fmt.Println("failed to get maps")
			return errIn
		}
	}

	var pokemon pokeApi.LocationPokemon
	pokemon, errIn = pokeApi.JsonToLocationPokemon(data)
	if errIn != nil {
		fmt.Printf("failed to convert json to location pokemon struct %v\n", errIn.Error())
		return errIn
	}

	fmt.Printf("Exploring %v...\n", area_name[0])
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokemon.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *configCommand, area_name []string) error {
	if len(area_name) != 1 {
		fmt.Print("Catch has only one arg (pokemon name)\n")
		return errors.New("too many args")
	}
	pokemonName := strings.ToLower(area_name[0])
	url := config.PokemonUrl + pokemonName
	data, err := cache.Get(url)
	var errIn error
	if !err {
		data, errIn = pokeApi.CallApiByUrl(url)
		if errIn != nil {
			fmt.Printf("failed to get pokemon %v\n", pokemonName)
			return errIn
		}
	}

	var pokemon pokeApi.Pokemon
	pokemon, errIn = pokeApi.JsonToPokemon(data)
	if errIn != nil {
		fmt.Printf("failed to convert json to pokemon data")
		return errIn
	}

	// Try to catch

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	catchChanges := pokemon.BaseExperience
	if rand.Intn(catchChanges) > (catchChanges - 30) {
		fmt.Printf("%v was caught!\n", pokemonName)
		config.Pokedex[pokemonName] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func getLocation(url string) (pokeApi.Location, error) {
	fmt.Println(url)
	data, err := cache.Get(url)
	var errIn error
	if !err {
		data, errIn = pokeApi.CallApiByUrl(url)
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
		"catch": {
			name:        "catch",
			description: "try to catch the named pokemon",
			callback:    commandCatch,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	baseLocationUrl := "https://pokeapi.co/api/v2/location/"

	config := configCommand{
		Pokedex:     make(map[string]pokeApi.Pokemon),
		PokemonUrl:  "https://pokeapi.co/api/v2/pokemon/",
		LocationUrl: "https://pokeapi.co/api/v2/location-area/",
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
