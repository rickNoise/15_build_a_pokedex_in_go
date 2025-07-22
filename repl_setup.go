package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokeapi"
	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokecache"
)

// cleanInput splits the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	textTrimmedToLower := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(textTrimmedToLower)
}

// initialise the repl environment for main.go
// returns an instance of config for the user and a scanner to read input
// also creates a cache to be used to minimise network calls
func ReplInitialisation() (*config, *bufio.Scanner) {
	locationCache, err := pokecache.NewCache(CACHE_LIFE_IN_SECONDS * time.Second)
	if err != nil {
		fmt.Print(fmt.Errorf("problem initialising cache in userConfig: %w", err))
	}
	var userConfig = &config{
		Next:          "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		Previous:      "",
		LocationCache: locationCache,
		Pokedex:       make(map[string]pokeapi.Pokemon),
	}
	scanner := bufio.NewScanner(os.Stdin)
	return userConfig, scanner
}

// config represents the user's state when exploring the Pokemon universe.
// Next and Previous are using to paginate through location areas.
// ExplorationIncrement sets how many new locations are shown when the user uses commands map or mapb. Basically the size of the "step" taken when exploring through the location space.
type config struct {
	Next          string
	Previous      string
	LocationCache *pokecache.Cache
	Pokedex       map[string]pokeapi.Pokemon // violating clean architecture
}

// cliCommand represents a command that can be called by the user from the CLI.
type cliCommand struct {
	name        string
	description string
	// callback takes userConfig and user prompt, split into words
	callback func(*config, []string) error
}

// GetCommands returns the hardcoded map of available commands
func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message.",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Explore the map. Displays the next 20 area names.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Explore back the way you came. Displays the previous 20 area names.",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area for Pokemon. e.g. \"explore <area name>\". Find area names by using \"map\" first.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon. e.g. \"catch <pokemon name>\". Use \"explore\" command to find Pokemon names.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "See details about a Pokemon. e.g. \"inspect <pokemon name>\". You must catch a Pokemon before you can inspect it.",
			callback:    commandInspect,
		},
	}
}
