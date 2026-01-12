package webserver

import (
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
			// unwrap to see if not found or not? maybe return sentinal error here
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Bad request"))
			return
		}

		// should I return this as ascii art or convert to img that I stream back? Probably the latter as, to begin
		// UI shouldn't need to format the ascii art
		w.WriteHeader(http.StatusOK)
		pokeimage.NewPokemonImage(image).Write(w)
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		return err
	}

	return nil
}
