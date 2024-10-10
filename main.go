package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	pokemonName := "snorlax"
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))
	panicIfErr(err)

	defer resp.Body.Close()
	var pokemonResponse PokemonResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&pokemonResponse)

	resp, err = http.Get(pokemonResponse.Sprites.Default)
	panicIfErr(err)
	defer resp.Body.Close()

	// read contents of image response
	fileName := fmt.Sprintf("./images/%s.jpg", pokemonName)
	file, fileErr := os.Create(fileName)
	panicIfErr(fileErr)

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	panicIfErr(err)
	fmt.Printf("image saved to %s", fileName)
}
