package commands

import "github.com/mirsafari/pokedexcli/internal/state"

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*state.DataStore, string) error
}

func Initialize() map[string]CliCommand {

	cmds := make(map[string]CliCommand)

	cmds["help"] = CliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    Help,
	}
	cmds["exit"] = CliCommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    Exit,
	}
	cmds["map"] = CliCommand{
		Name:        "map",
		Description: "Displays the names of 20 location areas in the Pokemon world",
		Callback:    Map,
	}
	cmds["mapb"] = CliCommand{
		Name:        "mapb",
		Description: "Similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations",
		Callback:    Mapb,
	}
	cmds["explore"] = CliCommand{
		Name:        "explore",
		Description: "See a list of all the Pokémon in a given area. Provide a location name as argument",
		Callback:    Explore,
	}
	cmds["catch"] = CliCommand{
		Name:        "catch",
		Description: "Catching Pokemon adds them to the user's Pokedex. Provide a Pokemon name as argument",
		Callback:    Catch,
	}
	cmds["inspect"] = CliCommand{
		Name:        "inspect",
		Description: "Prints the name, height, weight, stats and type(s) of the. Provide a Pokemon name as argument",
		Callback:    Inspect,
	}
	cmds["pokedex"] = CliCommand{
		Name:        "pokedex",
		Description: "Prints a list of all the names of the Pokemon the user has caught",
		Callback:    Pokedex,
	}
	return cmds
}
