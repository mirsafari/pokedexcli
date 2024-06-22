package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type LocationDetails struct {
	ID       int `json:"id"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationDetails(endpoint, location string, cache *pokecache.Cache) (LocationDetails, error) {
	locationDetails := LocationDetails{}

	data, inCache := cache.Get(location)
	if inCache == false {
		res, err := http.Get(endpoint + "location-area/" + location)
		if err != nil {
			return locationDetails, err
		}

		data, err = io.ReadAll(res.Body)
		res.Body.Close()

		if res.StatusCode == 404 {
			return locationDetails, errors.New(fmt.Sprintf("Location does not exist"))
		}

		if res.StatusCode > 299 {
			return locationDetails, errors.New(fmt.Sprintf("Response failed with status code %d", res.StatusCode))
		}

		cache.Add(location, data)
	}

	err := json.Unmarshal(data, &locationDetails)
	if err != nil {
		return locationDetails, errors.New("Error converting JSON response")
	}

	return locationDetails, nil
}
