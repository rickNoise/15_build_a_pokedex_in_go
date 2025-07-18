package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var commandMap map[string]cliCommand

func main() {
	var userConfig = config{
		Next:                 "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		Previous:             "",
		ExplorationIncrement: 20,
	}

	commandMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Explore the Pokemon world. Displays the next 20 locations.",
			callback:    commandExplore,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back the way you came. Displays the previous 20 locations shown by a \"Map\" command.",
			callback:    commandExploreBack,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for isRunning := true; isRunning; {
		var userPrompt []string
		fmt.Printf("Current user config:\n%+v\n", userConfig)
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("error parsing user input")
			log.Fatal(1)
		}
		userPrompt = cleanInput(scanner.Text())
		userPromptFirstWord := userPrompt[0]
		if userCommand, exists := commandMap[userPromptFirstWord]; exists {
			err := userCommand.callback(&userConfig)
			if err != nil {
				fmt.Print(fmt.Errorf("error running command: %w", err))
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
