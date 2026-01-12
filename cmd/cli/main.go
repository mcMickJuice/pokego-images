package main

import (
	"flag"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
)

// TODO implement multiple apps/clients:
// cli
// webserver - webserver could contain a redis cache to speed up pokemon fetch

func main() {
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	pokemonClient := pokemon.NewPokemonClient(*pokemonPtr)
	image, err := pokemonClient.GetPokemonSprite()

	if err != nil {
		panic(err)
	}

	pokeimage.NewPokemonImage(image).PrintAsciiArt()

}
