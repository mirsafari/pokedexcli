package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type PokemonDetails struct {
	Name      string `json:"name"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Height         int `json:"height"`
	Weight         int `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func GetPokemonDetails(pokemonName string, cache *pokecache.Cache) (PokemonDetails, error) {
	pokemonDetails := PokemonDetails{}
	data, inCache := cache.Get(pokemonName)

	if inCache == false {
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)
		if err != nil {
			return pokemonDetails, err
		}
		data, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode == 404 {
			return pokemonDetails, errors.New(fmt.Sprintf("Pokemon does not exist"))
		}
		if res.StatusCode > 299 {
			return pokemonDetails, errors.New(fmt.Sprintf("Response failed with status code %d", res.StatusCode))
		}

		cache.Add(pokemonName, data)
	}

	err := json.Unmarshal(data, &pokemonDetails)
	if err != nil {
		return pokemonDetails, errors.New("Error converting JSON response")
	}

	return pokemonDetails, nil
}
