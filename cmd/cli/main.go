package main

import (
	"context"
	"flag"
	"log"
	"mcmickjuice/pokego/internal/asciiimage"
	"mcmickjuice/pokego/internal/pokemon"
	"os"
)

func main() {
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	pokemonClient := pokemon.NewPokemonClient("https://pokeapi.co")
	ctx := context.Background()
	resp, err := pokemonClient.GetPokemon(ctx, *pokemonPtr)

	if err != nil {
		log.Fatalf("There was an error: %v", err)
	}

	image, err := pokemonClient.GetPokemonSpriteImage(context.Background(), resp)

	if err != nil {
		log.Fatalf("There was an error: %v", err)
	}

	err = asciiimage.NewAsciiImage(image).Write(os.Stdout)
	if err != nil {
		log.Printf("error writing output: %v", err)
	}
}
