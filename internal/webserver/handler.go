package webserver

import (
	"fmt"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
	"net/http"
)

func CreateWebserver() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/pokemon/{pokemon}", func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("pokemon")
		fmt.Printf("pokemon param: %s", param)

		pokemonClient := pokemon.NewPokemonClient(param)
		image, err := pokemonClient.GetPokemonSprite()

		if err != nil {
			// unwrap to see if not found or not? maybe return sentinal error here
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Bad request"))
			return
		}

		asciiArt := pokeimage.NewPokemonImage(image).AsciiArt()
		_, err = w.Write([]byte(asciiArt))
		if err != nil {
			fmt.Printf("error writing response: %v", err)
		}
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		return err
	}

	return nil
}
