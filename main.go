package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites PokemonSpritesResponse `json:"sprites"`
}

type PokemonSpritesResponse struct {
	Default string `json:"front_default"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/snorlax")
	panicIfErr(err)

	var pokemonResponse PokemonResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&pokemonResponse)

	fmt.Printf("Name: %s, Id: %d, Sprite: %s", pokemonResponse.Name, pokemonResponse.Id, pokemonResponse.Sprites.Default)
}
