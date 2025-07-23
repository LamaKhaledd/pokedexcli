package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/LamaKhaledd/pokedexcli/internal/pokeapi"
	"github.com/LamaKhaledd/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextLocationURL     *string
	previousLocationURL *string
}

func main() {
	pokeapi.Cache = pokecache.NewCache(5 * time.Minute)

	cfg := &config{}

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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Println("\nExiting Pokedex.")
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		cmd, found := commands[cmdName]

		if !found {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(cfg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
