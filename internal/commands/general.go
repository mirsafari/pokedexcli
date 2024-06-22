package commands

import (
	"fmt"
	"os"

	"github.com/mirsafari/pokedexcli/internal/state"
)

func Help(_ *state.DataStore, _ string) error {
	fmt.Println("Wellcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	// A bit ugly to allocate variable back again, but I don't know the workaround
	inputs := Initialize()

	for _, v := range inputs {
		fmt.Println(v.Name, ":", v.Description)
	}

	fmt.Println("")
	return nil
}

func Exit(_ *state.DataStore, _ string) error {
	os.Exit(0)
	return nil
}
