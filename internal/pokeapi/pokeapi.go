package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/LamaKhaledd/pokedexcli/internal/pokecache"
)

type LocationArea struct {
	Name string `json:"name"`
}

type LocationAreasResponse struct {
	Results  []LocationArea `json:"results"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
}

var Cache *pokecache.Cache

func GetLocationAreas(url string) ([]string, *string, *string, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	// Check cache first
	if data, found := Cache.Get(url); found {
		fmt.Println("📦 Loaded from cache:", url)
		var parsed LocationAreasResponse
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, nil, nil, err
		}

		var names []string
		for _, area := range parsed.Results {
			names = append(names, area.Name)
		}
		return names, parsed.Next, parsed.Previous, nil
	}

	// Fetch from network
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, errors.New("failed to fetch location areas")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, err
	}

	// Cache the result
	Cache.Add(url, body)

	var data LocationAreasResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, nil, nil, err
	}

	var names []string
	for _, area := range data.Results {
		names = append(names, area.Name)
	}

	return names, data.Next, data.Previous, nil
}

func GetPokemonInLocationArea(name string) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)

	// Check cache first
	if data, found := Cache.Get(url); found {
		fmt.Println("📦 Loaded from cache:", url)
		return parsePokemonFromLocationData(data)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch location area data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Cache it
	Cache.Add(url, body)

	return parsePokemonFromLocationData(body)
}

func parsePokemonFromLocationData(data []byte) ([]string, error) {
	var resp struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	var names []string
	for _, encounter := range resp.PokemonEncounters {
		names = append(names, encounter.Pokemon.Name)
	}

	return names, nil
}
