package state

import (
	"time"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/pokecache"
)

type DataStore struct {
	APIEndpoint         string
	NextLocationURL     string
	PreviousLocationURL string
	APIResponseCache    *pokecache.Cache
	PokemonContainer    map[string]pokeapi.PokemonDetails
}

func Initialize(apiEndpoint string, cachePurgeInterval time.Duration) DataStore {

	locationEndpoint := apiEndpoint + "location-area?offset=0&limit=20"

	return DataStore{
		APIEndpoint:      apiEndpoint,
		NextLocationURL:  locationEndpoint,
		APIResponseCache: pokecache.NewCache(cachePurgeInterval),
		PokemonContainer: make(map[string]pokeapi.PokemonDetails),
	}
}
