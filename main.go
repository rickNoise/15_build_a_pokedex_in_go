package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var commandMap map[string]cliCommand

func main() {
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
	}

	scanner := bufio.NewScanner(os.Stdin)

	for isRunning := true; isRunning; {
		var userPrompt []string
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println("error parsing user input")
			log.Fatal(1)
		}
		userPrompt = cleanInput(scanner.Text())
		userPromptFirstWord := userPrompt[0]
		if userCommand, exists := commandMap[userPromptFirstWord]; exists {
			userCommand.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}
