package main

import (
	"flag"
	"log"
	"mcmickjuice/pokego/internal/asciiimage"
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

	err = asciiimage.NewAsciiImage(image).Write(os.Stdout)
	if err != nil {
		log.Printf("error writing output: %v", err)
	}
}
