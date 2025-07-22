package main

import (
	"fmt"
	"log"
)

/* CONSTANTS */
const CACHE_LIFE_IN_SECONDS = 60

func main() {
	// initalise repl environment
	userConfig, scanner := ReplInitialisation()

	// show help on start
	GetCommands()["help"].callback(userConfig, nil)

	// cli user input loop
	for isRunning := true; isRunning; {
		fmt.Printf("\nPokedex > ")
		if !scanner.Scan() {
			fmt.Println("error parsing user input")
			log.Fatal(1)
		}

		userPrompt := cleanInput(scanner.Text())
		if len(userPrompt) == 0 {
			continue
		}

		userCommand := userPrompt[0]
		if userCommand, exists := GetCommands()[userCommand]; exists {
			err := userCommand.callback(userConfig, userPrompt)
			if err != nil {
				fmt.Print(fmt.Errorf("error running command: %w", err))
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
