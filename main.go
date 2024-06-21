package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type config struct {
	Url      string
	Next     string
	Previous string
	Cache    *pokecache.Cache
}
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var (
	commands        = make(map[string]cliCommand)
	cliName  string = "Pokedex"
	cfg      config = config{
		Url:      "",
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0",
		Previous: "",
	}
)

func initializeCommands() {
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the names of 20 location areas in the Pokemon world",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations",
		callback:    commandMapb,
	}
	return
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func commandHelp(*config) error {
	fmt.Println("Wellcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, v := range commands {
		fmt.Println(v.name, ":", v.description)
	}

	fmt.Println("")
	return nil
}

func commandExit(*config) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	locations, err := pokeapi.GetLocations(cfg.Next, cfg.Cache)

	if err != nil {
		return err
	}

	printLocations(locations)
	cfg.Previous = locations.Previous
	cfg.Next = locations.Next

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == "" {
		return errors.New("Can not go backwards on first location")
	}
	locations, err := pokeapi.GetLocations(cfg.Previous, cfg.Cache)

	if err != nil {
		return err
	}

	printLocations(locations)
	cfg.Previous = locations.Previous
	cfg.Next = locations.Next

	return nil
}

func printLocations(locations pokeapi.Locations) {

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return
}

func main() {
	initializeCommands()
	cfg.Cache = pokecache.NewCache(time.Minute * 2)
	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		text := reader.Text()
		if command, exists := commands[text]; exists {
			err := command.callback(&cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
		printPrompt()
	}
}
