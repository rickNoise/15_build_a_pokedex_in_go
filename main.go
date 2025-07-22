package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokeapi"
	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokecache"
)

/* CONSTANTS */
const CACHE_LIFE_IN_SECONDS = 60

func main() {
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

	// show help on start
	GetCommands()["help"].callback(userConfig, nil)

	// cli loop
	for isRunning := true; isRunning; {
		var userPrompt []string
		fmt.Printf("\nPokedex > ")
		if !scanner.Scan() {
			fmt.Println("error parsing user input")
			log.Fatal(1)
		}
		userPrompt = cleanInput(scanner.Text())
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
