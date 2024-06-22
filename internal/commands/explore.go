package commands

import (
	"fmt"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/state"
)

func Explore(data *state.DataStore, locationName string) error {

	fmt.Println("Exploring " + locationName + "...")
	locationDetails, err := pokeapi.GetLocationDetails(data.APIEndpoint, locationName, data.APIResponseCache)

	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationDetails.PokemonEncounters {
		fmt.Println(" -", pokemon.Pokemon.Name)
	}

	return nil
}
