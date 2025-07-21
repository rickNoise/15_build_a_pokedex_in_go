package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rickNoise/15_build_a_pokedex_in_go/internal/pokecache"
)

// struct to capture json response from GetLocationAreas
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
func GetLocationAreas(url string, locationCache *pokecache.Cache) ([]string, string, string, error) {
	// if passed URL is empty, we haven't explored at all yet
	// kick start with a search with 0 offset
	if url == "" {
		return []string{}, "", "", errors.New("error: empty url string provided")
	}

	results, foundInCache := locationCache.Get(url)
	if !foundInCache {
		fmt.Println("Not found in cache, calling API...")
		res, err := http.Get(url)
		if err != nil {
			return []string{}, "", "", errors.New("error: Could not GET Location Areas")
		}
		defer res.Body.Close()

		results, err = io.ReadAll(res.Body)
		if err != nil {
			return []string{}, "", "", errors.New("error: could not read response body")
		}

		locationCache.Add(url, results)
	} else {
		fmt.Println("Using cache on GetLocationAreas...")
	}

	var LocationAreas LocationAreasResponse
	err := json.Unmarshal(results, &LocationAreas)
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

type LocationAreaSpecificsResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		Version struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
		MaxChance        int `json:"max_chance"`
		EncounterDetails []struct {
			MinLevel        int           `json:"min_level"`
			MaxLevel        int           `json:"max_level"`
			ConditionValues []interface{} `json:"condition_values"`
			Chance          int           `json:"chance"`
			Method          struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
		} `json:"encounter_details"`
	} `json:"version_details"`
}

func GetPokemonInArea(areaName string, locationCache *pokecache.Cache) ([]PokemonEncounter, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	results, foundInCache := locationCache.Get(url)
	if !foundInCache {
		fmt.Println("Not found in cache, calling API...")
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error: Could not get details for area %v", areaName)
		}
		defer res.Body.Close()

		results, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, errors.New("error: could not read response body")
		}

		locationCache.Add(url, results)
	} else {
		fmt.Println("Using cache...")
	}

	var LocationAreaSpecificResponseData LocationAreaSpecificsResponse
	err := json.Unmarshal(results, &LocationAreaSpecificResponseData)
	if err != nil {
		return nil, errors.New("error: could not Unmarshall results from res Reader")
	}

	var PokemonEncounters []PokemonEncounter
	for _, item := range LocationAreaSpecificResponseData.PokemonEncounters {
		PokemonEncounters = append(PokemonEncounters, item)
	}

	return PokemonEncounters, nil
}
