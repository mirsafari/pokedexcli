package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationDetails struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocations(url string, cache *pokecache.Cache) (Locations, error) {

	locations := Locations{}
	data, inCache := cache.Get(url)

	if inCache == false {
		res, err := http.Get(url)
		if err != nil {
			return locations, err
		}
		data, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return locations, errors.New(fmt.Sprintf("Response failed with status code %d", res.StatusCode))
		}

		cache.Add(url, data)
	} else {
		fmt.Println("SERVED FROM CACHE")
	}

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return locations, errors.New("Error converting JSON response")
	}

	return locations, nil
}

func GetLocationDetails(location string, cache *pokecache.Cache) (LocationDetails, error) {

	locationDetails := LocationDetails{}
	data, inCache := cache.Get(location)

	if inCache == false {
		res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + location)
		if err != nil {
			return locationDetails, err
		}
		data, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return locationDetails, errors.New(fmt.Sprintf("Response failed with status code %d", res.StatusCode))
		}

		cache.Add(location, data)
	} else {
		fmt.Println("SERVED FROM CACHE")
	}

	err := json.Unmarshal(data, &locationDetails)
	if err != nil {
		return locationDetails, errors.New("Error converting JSON response")
	}

	return locationDetails, nil
}
