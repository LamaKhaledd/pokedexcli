package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/LamaKhaledd/pokedexcli/internal/pokeapi"
)

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	words := strings.Fields(lowered)
	return words
}

func commandHelp(cmds map[string]cliCommand) func(*config) error {
	return func(cfg *config) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		for _, cmd := range cmds {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		fmt.Println()
		return nil
	}
}

func commandExit() func(*config) error {
	return func(cfg *config) error {
		fmt.Println("Closing the Pokedex... Goodbye!")
		os.Exit(0)
		return nil
	}
}

func commandMap(cfg *config) error {
	url := ""
	if cfg.nextLocationURL != nil {
		url = *cfg.nextLocationURL
	}

	locations, next, prev, err := pokeapi.GetLocationAreas(url)
	if err != nil {
		return err
	}

	for _, loc := range locations {
		fmt.Println(loc)
	}

	cfg.nextLocationURL = next
	cfg.previousLocationURL = prev

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.previousLocationURL == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	locations, next, prev, err := pokeapi.GetLocationAreas(*cfg.previousLocationURL)
	if err != nil {
		return err
	}

	for _, loc := range locations {
		fmt.Println(loc)
	}

	cfg.nextLocationURL = next
	cfg.previousLocationURL = prev

	return nil
}
