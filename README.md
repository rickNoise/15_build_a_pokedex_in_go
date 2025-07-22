# Build a Pokedex in Go

HTTP client project for boot.dev.

# How to Run

1. Clone repo.
2. Type "go run ." from the root dir in your terminal and press Enter.
3. Interact with the cli tool.

# Example Usage

The first word you enter is interpreted as a command.
Some commands use the next word as a cli argument for the command.
Type "help" to see available commands.

1. "map" shows next 20 areas names. Use this to see a list of areas.
2. "explore <area name>" using an area name found by using "map". Shows a list of Pokemon in the area.
3. "catch <pokemon name>" using a name found by exploring an area. More advanced Pokemon are less likely to be caught on the first attempt.
4. "inspect <pokemon name>" shows details of a caught Pokemon. You can only inspect Pokemon you've already caught.

# Implementation Details

- Calls [PokeAPI](https://pokeapi.co/docs/v2) for data.
- Uses local 60-second cache to reduce API calls.
- Unmarshals JSON responses from API into Go structs. [JSON to GO](https://transform.tools/json-to-go) was very useful for achieving this.
