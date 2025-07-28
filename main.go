package main

import (
	"fmt"
	"os"
	"time"
	"github.com/chzyer/readline"
	"github.com/LamaKhaledd/pokedexcli/internal/pokeapi"
	"github.com/LamaKhaledd/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	nextLocationURL     *string
	previousLocationURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
}


func main() {
	pokeapi.Cache = pokecache.NewCache(5 * time.Minute)

	cfg := &config{
	caughtPokemon: make(map[string]pokeapi.Pokemon),
	}


	commands := make(map[string]cliCommand)

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit(),
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(commands),
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "View next 20 location areas",
		callback:    commandMap,
	}

	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "View previous 20 location areas",
		callback:    commandMapBack,
	}

	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explore a location area and see Pokémon there",
		callback:    commandExplore,
	}

	commands["catch"] = cliCommand{
	name:        "catch",
	description: "Try to catch a Pokemon by name",
	callback:    commandCatch,
	}


	rl, err := readline.New("Pokedex > ")
	if err != nil {
		fmt.Println("Failed to start input reader:", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("\nExiting Pokedex.")
			break
		}

		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		cmd, found := commands[cmdName]

		if !found {
			fmt.Println("Unknown command")
			continue
		}

		args := words[1:]
		err = cmd.callback(cfg, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}