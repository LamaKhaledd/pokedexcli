package main

import (
	"fmt"
	"os"
	"strings"
	"math/rand"
	"time"

	"github.com/LamaKhaledd/pokedexcli/internal/pokeapi"
)

type Stat struct {
	BaseStat int `json:"base_stat"`
	StatInfo struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Type struct {
	TypeInfo struct {
		Name string `json:"name"`
	} `json:"type"`
}

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	words := strings.Fields(lowered)
	return words
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range commandsMap {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config, args []string) error {
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

func commandMapBack(cfg *config, args []string) error {
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

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: explore <location_area_name>")
		return nil
	}

	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)

	pokemonNames, err := pokeapi.GetPokemonInLocationArea(areaName)
	if err != nil {
		return err
	}

	if len(pokemonNames) == 0 {
		fmt.Println("No Pokémon found in this location area.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, name := range pokemonNames {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: catch <pokemon_name>")
		return nil
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemon: %w", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	chance := 100 - pokemon.BaseExperience
	if chance < 10 {
		chance = 10
	}

	roll := r.Intn(100)
	if roll < chance {
		cfg.caughtPokemon[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: inspect <pokemon_name>")
		return nil
	}

	name := args[0]

	p, caught := cfg.caughtPokemon[name]
	if !caught {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf("  -%s: %d\n", stat.StatInfo.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.TypeInfo.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("You haven't caught any Pokemon yet.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
