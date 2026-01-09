package main

import (
	"flag"
	"mcmickjuice/pokego/pokeimage"
	"mcmickjuice/pokego/pokemon"
)

type PokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites PokemonSpritesResponse `json:"sprites"`
}

type PokemonSpritesResponse struct {
	Default     string `json:"front_default"`
	BackDefault string `json:"back_default"`
	FrontShiny  string `json:"front_shiny"`
}

type AllPokemonsResponse struct {
	Results []AllPokemonResponse `json:"results"`
}

type AllPokemonResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
const POKEMON_DETAIL_URL = "https://pokeapi.co/api/v2/pokemon/%s"

func main() {
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	pokemonClient := pokemon.NewPokemonClient(*pokemonPtr)
	image, err := pokemonClient.GetPokemonSprite()

	if err != nil {
		panic(err)
	}

	pokeimage.NewPokemonImage(image).ToAsciiArt()

}
