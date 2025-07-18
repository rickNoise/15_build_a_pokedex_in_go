package main

import (
	"fmt"
	"os"
	"strings"
)

// cleanInput splits the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func cleanInput(text string) []string {
	textTrimmedToLower := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(textTrimmedToLower)
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	welcomeLine := "Welcome to the Pokedex!"
	usageLine := "Usage:\n"

	fmt.Println(welcomeLine)
	fmt.Println(usageLine)
	for _, value := range commandMap {
		fmt.Println(value.name + ":", value.description)
	}
	return nil
}
