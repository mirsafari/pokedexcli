# pokedexcli
Pokedex-like REPL in Go that uses the PokeAPI to fetch data about Pokemon.

This is a tutorial task on [boot.dev](https://boot.dev).

## Usage
```bash
go build && ./pokedexcli

Pokedex> help
Wellcome to the Pokedex!
Usage:

map : Displays the names of 20 location areas in the Pokemon world
mapb : Similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations
explore : See a list of all the Pokémon in a given area. Provide a location name as argument
catch : Catching Pokemon adds them to the user's Pokedex. Provide a Pokemon name as argument
inspect : Prints the name, height, weight, stats and type(s) of the. Provide a Pokemon name as argument
pokedex : Prints a list of all the names of the Pokemon the user has caught
help : Displays a help message
exit : Exit the Pokedex
```
