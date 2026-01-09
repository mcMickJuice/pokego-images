package main

import (
	"flag"
	"mcmickjuice/pokego/pokeimage"
	"mcmickjuice/pokego/pokemon"
)

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
