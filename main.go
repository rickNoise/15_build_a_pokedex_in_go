package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokecache"
)

var commandMapper map[string]cliCommand

func main() {
	locationCache, err := pokecache.NewCache(5 * time.Second)
	if err != nil {
		fmt.Print(fmt.Errorf("probably initialising cache in userConfig: %w", err))
	}
	var userConfig = &config{
		Next:          "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		Previous:      "",
		LocationCache: locationCache,
	}

	commandMapper = map[string]cliCommand{
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
	}

	scanner := bufio.NewScanner(os.Stdin)

	commandMapper["help"].callback(userConfig, nil)
	for isRunning := true; isRunning; {
		var userPrompt []string
		// fmt.Printf("Current user config:\n%+v\n", userConfig)
		fmt.Printf("\nPokedex > ")
		if !scanner.Scan() {
			fmt.Println("error parsing user input")
			log.Fatal(1)
		}
		userPrompt = cleanInput(scanner.Text())
		userPromptFirstWord := userPrompt[0]
		if userCommand, exists := commandMapper[userPromptFirstWord]; exists {
			err := userCommand.callback(userConfig, userPrompt)
			if err != nil {
				fmt.Print(fmt.Errorf("error running command: %w", err))
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
