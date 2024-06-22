package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/mirsafari/pokedexcli/internal/pokeapi"
	"github.com/mirsafari/pokedexcli/internal/state"
)

func Catch(data *state.DataStore, pokemonName string) error {

	pokemonDetails, err := pokeapi.GetPokemonDetails(data.APIEndpoint, pokemonName, data.APIResponseCache)

	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")

	randomValue := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
	catchProbability := 1 / (1 + float64(pokemonDetails.BaseExperience)/100.0)

	if randomValue < catchProbability {
		fmt.Println(pokemonName + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		data.PokemonContainer[pokemonName] = pokemonDetails
		return nil
	}

	return errors.New(pokemonName + " escaped!")
}

func Inspect(data *state.DataStore, pokemonName string) error {

	pokemon, caught := data.PokemonContainer[pokemonName]
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

func Pokedex(data *state.DataStore, _ string) error {
	fmt.Println("Your Pokedex:")
	for k := range data.PokemonContainer {
		fmt.Println(" -", k)
	}
	return nil
}
