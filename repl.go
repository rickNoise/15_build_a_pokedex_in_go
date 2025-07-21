package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokeapi"
	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokecache"
)

// cleanInput splits the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	textTrimmedToLower := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(textTrimmedToLower)
}

// config represents the user's state when exploring the Pokemon universe.
// Next and Previous are using to paginate through location areas.
// ExplorationIncrement sets how many new locations are shown when the user uses commands map or mapb. Basically the size of the "step" taken when exploring through the location space.
type config struct {
	Next          string
	Previous      string
	LocationCache *pokecache.Cache
}

// cliCommand represents a command that can be called by the user from the CLI.
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(userConfig *config) error {
	userConfig.LocationCache.Stop()
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(userConfig *config) error {
	welcomeLine := "Welcome to the Pokedex!"
	usageLine := "Usage:\n"

	fmt.Println(welcomeLine)
	fmt.Println(usageLine)
	for _, value := range commandMap {
		fmt.Println(value.name+":", value.description)
	}
	return nil
}

func commandExplore(userConfig *config) error {
	locationSlice, nextURL, prevURL, err := pokeapi.GetLocationAreas(
		userConfig.Next,
		userConfig.LocationCache,
	)
	if err != nil {
		return fmt.Errorf("error: map command failed: %w", err)
	}

	userConfig.Previous = prevURL
	userConfig.Next = nextURL

	for _, location := range locationSlice {
		fmt.Println(location)
	}

	return nil
}

func commandExploreBack(userConfig *config) error {
	// check to see if user is at the beginning of the exploration map.
	if userConfig.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationSlice, nextURL, prevURL, err := pokeapi.GetLocationAreas(
		userConfig.Previous,
		userConfig.LocationCache,
	)
	if err != nil {
		return fmt.Errorf("error: mapb command failed: %w", err)
	}

	userConfig.Next = nextURL
	userConfig.Previous = prevURL

	for _, location := range locationSlice {
		fmt.Println(location)
	}

	return nil
}
