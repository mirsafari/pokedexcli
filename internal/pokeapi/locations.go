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
	}

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return locations, errors.New("Error converting JSON response")
	}

	return locations, nil
}
