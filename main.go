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

type AllPokemonsResponse struct {
	Results []AllPokemonResponse `json:"results"`
}

type AllPokemonResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
const POKEMON_DETAIL_URL = "https://pokeapi.co/api/v2/pokemon/%s"

func main() {
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, "pikachu"))
	panicIfErr(err)

	var pokemonResponse PokemonResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&pokemonResponse)

	fmt.Printf("Name: %s, Id: %d, Sprite: %s", pokemonResponse.Name, pokemonResponse.Id, pokemonResponse.Sprites.Default)
}
