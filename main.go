package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")

        if !scanner.Scan() {
            fmt.Println("\nExiting Pokedex.")
            break
        }

        input := scanner.Text()
        words := cleanInput(input)

        if len(words) > 0 {
            fmt.Printf("Your command was: %s\n", words[0])
        }
    }
}
