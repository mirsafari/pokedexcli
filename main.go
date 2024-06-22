package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mirsafari/pokedexcli/internal/commands"
	"github.com/mirsafari/pokedexcli/internal/state"
)

const (
	cliName               string        = "Pokedex"
	APIEndpoint           string        = "https://pokeapi.co/api/v2/"
	APICachePurgeInterval time.Duration = time.Minute * 2
)

func main() {
	inputOptions := commands.Initialize()
	data := state.Initialize(APIEndpoint, APICachePurgeInterval)
	reader := bufio.NewScanner(os.Stdin)

	printPrompt()

	for reader.Scan() {
		command, argument := parseCommad(reader.Text())

		if execute, exists := inputOptions[command]; exists {
			err := execute.Callback(&data, argument)

			if err != nil {
				fmt.Println(err)
			}

		} else {
			fmt.Println("Command not found. Type help to list available commands")
		}

		printPrompt()
	}
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func parseCommad(line string) (string, string) {
	input := strings.Fields(line)

	cmd := input[0]
	arg := ""

	if len(input) == 2 {
		arg = input[1]
	}
	return cmd, arg
}
