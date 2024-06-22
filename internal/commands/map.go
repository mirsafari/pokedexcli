package commands

import (
	"errors"
	"fmt"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/state"
)

func Map(data *state.DataStore, _ string) error {
	locations, err := pokeapi.GetLocations(data.NextLocationURL, data.APIResponseCache)

	if err != nil {
		return err
	}

	printMap(locations)
	updateLocationURL(data, locations.Previous, locations.Next)

	return nil
}

func Mapb(data *state.DataStore, _ string) error {
	if data.PreviousLocationURL == "" {
		return errors.New("Can not go backwards on first location")
	}
	locations, err := pokeapi.GetLocations(data.PreviousLocationURL, data.APIResponseCache)

	if err != nil {
		return err
	}

	printMap(locations)
	updateLocationURL(data, locations.Previous, locations.Next)

	return nil
}

func printMap(locations pokeapi.Locations) {
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
}

func updateLocationURL(data *state.DataStore, previous, next string) {
	data.PreviousLocationURL = previous
	data.NextLocationURL = next

	return
}
