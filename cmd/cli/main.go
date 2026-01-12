package main

import (
	"flag"
	"fmt"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
)

type PrintWriter struct {
}

func (p PrintWriter) Write(input []byte) (int, error) {
	fmt.Println(string(input))
	return 0, nil
}

func main() {
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	pokemonClient := pokemon.NewPokemonClient(*pokemonPtr)
	image, err := pokemonClient.GetPokemonSprite()

	if err != nil {
		panic(err)
	}

	pokeimage.NewPokemonImage(image).Write(PrintWriter{})

}
