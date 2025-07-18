package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// PokePrint is a function for testing the internal package import is working.
func PokePrint() {
	fmt.Println("pokeapi internal package test print...")
}

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// GetLocationAreas pulls down number of locations areas from the API using the given "url".
// Returns a slice of location area names, the Next url for the next page of results, and an error.
func GetLocationAreas(url string, increment int) ([]string, string, string, error) {
	// if passed URL is empty, we haven't explored at all yet
	// kick start with a search with 0 offset
	if url == "" {
		url = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=%d&offset=0", increment)
	}

	res, err := http.Get(url)
	if err != nil {
		return []string{}, "", "", errors.New("error: Could not GET Location Areas")
	}
	defer res.Body.Close()

	results, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, "", "", errors.New("error: could not read response body")
	}

	var LocationAreas LocationAreasResponse
	err = json.Unmarshal(results, &LocationAreas)
	if err != nil {
		return []string{}, "", "", errors.New("error: could not Unmarshall results from res Reader")
	}

	var LocationAreaNames []string
	for _, result := range LocationAreas.Results {
		LocationAreaNames = append(LocationAreaNames, result.Name)
	}

	nextURL := LocationAreas.Next
	prevURL := ""
	if LocationAreas.Previous != nil {
		prevURL = *LocationAreas.Previous
	}

	return LocationAreaNames, nextURL, prevURL, nil
}
