package main

import (
	"flag"
	"log"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
	"os"
)

func main() {
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	pokemonClient := pokemon.NewPokemonClient(*pokemonPtr)
	image, err := pokemonClient.GetPokemonSprite()

	if err != nil {
		log.Fatalf("There was an error: %v", err)
	}

	pokeimage.NewPokemonImage(image).Write(os.Stdout)
}
