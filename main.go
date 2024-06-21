package main

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type config struct {
	Url              string
	Next             string
	Previous         string
	Cache            *pokecache.Cache
	PokemonContainer map[string]pokeapi.PokemonDetails
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
		description: "See a list of all the Pokémon in a given area. Provide a location name as argument",
		callback:    commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Catching Pokemon adds them to the user's Pokedex. Provide a Pokemon name as argument",
		callback:    commandCatch,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Prints the name, height, weight, stats and type(s) of the. Provide a Pokemon name as argument",
		callback:    commandInspect,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "Prints a list of all the names of the Pokemon the user has caught",
		callback:    commandPokedex,
	}
	return
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func commandHelp(cfg *config, _ string) error {
	fmt.Println("Wellcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, v := range commands {
		fmt.Println(v.name, ":", v.description)
	}

	fmt.Println("")
	return nil
}

func commandExit(cfg *config, _ string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, _ string) error {
	locations, err := pokeapi.GetLocations(cfg.Next, cfg.Cache)

	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	cfg.Previous = locations.Previous
	cfg.Next = locations.Next

	return nil
}

func commandMapb(cfg *config, _ string) error {
	if cfg.Previous == "" {
		return errors.New("Can not go backwards on first location")
	}
	locations, err := pokeapi.GetLocations(cfg.Previous, cfg.Cache)

	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

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

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationDetails.PokemonEncounters {
		fmt.Println(" -", pokemon.Pokemon.Name)
	}

	return nil

}

func commandCatch(cfg *config, pokemonName string) error {

	pokemonDetails, err := pokeapi.GetPokemonDetails(pokemonName, cfg.Cache)

	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")

	randomValue := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
	catchProbability := 1 / (1 + float64(pokemonDetails.BaseExperience)/100.0)

	if randomValue < catchProbability {
		fmt.Println(pokemonName + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.PokemonContainer[pokemonName] = pokemonDetails
		return nil
	}

	return errors.New(pokemonName + " escaped!")
}

func commandInspect(cfg *config, pokemonName string) error {

	pokemon, caught := cfg.PokemonContainer[pokemonName]
	if !caught {
		return errors.New("you have not caught that pokemon")
	}

	const tpl = `Name: {{.Name}}
Height: {{.Height}}
Weight: {{.Weight}}
Stats:
{{- range .Stats }}
  -{{ .Stat.Name }}: {{ .BaseStat }}
{{- end }}
Types:
{{- range .Types }}
  - {{ .Type.Name }}
{{- end }}
`
	tmpl, err := template.New("pokemon").Parse(tpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, pokemon)
	if err != nil {
		return err
	}

	return nil
}

func commandPokedex(cfg *config, _ string) error {
	fmt.Println("Your Pokedex:")
	for k := range cfg.PokemonContainer {
		fmt.Println(" -", k)
	}
	return nil
}

func main() {
	initializeCommands()
	cfg.Cache = pokecache.NewCache(time.Minute * 2)
	cfg.PokemonContainer = make(map[string]pokeapi.PokemonDetails)
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
