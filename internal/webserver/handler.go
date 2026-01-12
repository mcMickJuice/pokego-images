package webserver

import (
	"errors"
	"fmt"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
	"net/http"
)

func CreateWebserver() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/pokemon/{pokemon}", func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("pokemon")

		pokemonClient := pokemon.NewPokemonClient(param)
		image, err := pokemonClient.GetPokemonSprite()

		if err != nil {
			if errors.Is(err, pokemon.ErrPokemonNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, "Pokemon %s Not Found", param)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "Unknown error")
			return
		}

		w.WriteHeader(http.StatusOK)
		err = pokeimage.NewPokemonImage(image).Write(w)
		if err != nil {
			fmt.Printf("error writing to response: %v", err)
		}
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		return err
	}

	return nil
}
