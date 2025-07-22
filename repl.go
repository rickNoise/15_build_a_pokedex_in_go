package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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

func commandExit(userConfig *config, userPrompt []string) error {
	userConfig.LocationCache.Stop()
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(userConfig *config, userPrompt []string) error {
	welcomeLine := "Welcome to the Pokedex!"
	usageLine := "Usage:\n"

	fmt.Println(welcomeLine)
	fmt.Println(usageLine)
	for _, value := range GetCommands() {
		fmt.Println(value.name+":", value.description)
	}
	return nil
}

func commandMap(userConfig *config, userPrompt []string) error {
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

func commandMapBack(userConfig *config, userPrompt []string) error {
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

func commandExplore(userConfig *config, userPrompt []string) error {
	if len(userPrompt) < 2 {
		return errors.New("you must provide an area to explore after the \"explore\" command.\nE.g. \"explore <area name>\"")
	}
	userProvidedAreaName := userPrompt[1]

	pokemonInAreaSlice, err := pokeapi.GetPokemonInArea(userProvidedAreaName, userConfig.LocationCache)
	if err != nil {
		return errors.New("error: problem getting Pokemon in area")
	}

	for _, pokemon := range pokemonInAreaSlice {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(userConfig *config, userPrompt []string) error {
	if len(userPrompt) < 2 {
		return errors.New("you must provide an pokemon name after the \"catch\" command")
	}
	userProvidedPokemonName := userPrompt[1]

	if _, ok := userConfig.Pokedex[userProvidedPokemonName]; ok {
		return fmt.Errorf("you already have %s in your Pokedex", userProvidedPokemonName)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", userProvidedPokemonName)

	PokemonDetails, err := pokeapi.GetPokemonDetails(userProvidedPokemonName, userConfig.LocationCache)
	if err != nil {
		return errors.New("error: problem getting Pokemon details")
	}

	// The base experience gained for defeating this Pokémon (int).
	pokemonBaseExperience := PokemonDetails.BaseExperience

	// logic for determining if catch attempt is successful
	baseExpCapped := min(pokemonBaseExperience, 400)
	fmt.Printf("Base Experience: %d\n", baseExpCapped)

	randChance := 30 * (rand.Intn(9) + 1)
	fmt.Printf("randChance: %d\n", randChance)

	if randChance > baseExpCapped {
		fmt.Println(":) Pokemon caught!")
		userConfig.Pokedex[userProvidedPokemonName] = PokemonDetails
	} else {
		fmt.Println(";( Pokemon got away!")
	}

	return nil
}

func commandInspect(userConfig *config, userPrompt []string) error {
	if len(userPrompt) < 2 {
		return errors.New("you must provide an pokemon name after the \"inspect\" command")
	}
	userProvidedPokemonName := userPrompt[1]

	p, ok := userConfig.Pokedex[userProvidedPokemonName]
	if !ok {
		fmt.Printf("%s is not in your Pokedex. You must catch a Pokemon before you can inspect it.\n", userProvidedPokemonName)
		return nil
	}
	fmt.Println("Name:", p.Name)
	fmt.Println("Height:", p.Height)
	fmt.Println("Weight:", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("-%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Println("-", t.Type.Name)
	}
	return nil
}
