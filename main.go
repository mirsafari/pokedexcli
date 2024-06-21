package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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
	callback    func(*config, string) error
}

var (
	commands        = make(map[string]cliCommand)
	cliName  string = "Pokedex"
	cfg      config = config{
		Url:      "",
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
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
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "See a list of all the Pokémon in a given area. Provide a location name",
		callback:    commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Catching Pokemon adds them to the user's Pokedex. Provide a Pokemon name",
		callback:    commandExplore,
	}
	return
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func commandHelp(cfg *config, arg string) error {
	fmt.Println("Wellcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, v := range commands {
		fmt.Println(v.name, ":", v.description)
	}

	fmt.Println("")
	return nil
}

func commandExit(cfg *config, arg string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, arg string) error {
	locations, err := pokeapi.GetLocations(cfg.Next, cfg.Cache)

	if err != nil {
		return err
	}

	printLocations(locations)
	cfg.Previous = locations.Previous
	cfg.Next = locations.Next

	return nil
}

func commandMapb(cfg *config, arg string) error {
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

func commandExplore(cfg *config, location string) error {

	fmt.Println("Exploring " + location + "...")
	locationDetails, err := pokeapi.GetLocationDetails(location, cfg.Cache)

	if err != nil {
		return err
	}

	printPokemons(locationDetails)

	return nil

}

func commandCatch(cfg *config, pokemon string) error {

	return nil

}
func printLocations(locations pokeapi.Locations) {

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return
}

func printPokemons(location pokeapi.LocationDetails) {

	fmt.Println("Found Pokemon:")
	for _, pokemon := range location.PokemonEncounters {
		fmt.Println(" -", pokemon.Pokemon.Name)
	}

	return
}

func main() {
	initializeCommands()
	cfg.Cache = pokecache.NewCache(time.Minute * 2)
	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		input := strings.Fields(reader.Text())
		cmd := input[0]
		arg := ""
		if len(input) == 2 {
			arg = input[1]
		}

		if command, exists := commands[cmd]; exists {
			err := command.callback(&cfg, arg)
			if err != nil {
				fmt.Println(err)
			}
		}
		printPrompt()
	}
}
