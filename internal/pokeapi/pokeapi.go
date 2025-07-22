package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
}

type LocationAreasResponse struct {
	Results  []LocationArea `json:"results"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
}

func GetLocationAreas(url string) ([]string, *string, *string, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

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
