package webserver

import (
	"errors"
	"fmt"
	"log"
	"mcmickjuice/pokego/internal/pokeimage"
	"mcmickjuice/pokego/internal/pokemon"
	"net/http"
)

type PokemonWebServer struct {
	addr string
}

func NewPokemonWebServer(addr string) PokemonWebServer {
	return PokemonWebServer{addr}
}

func (s PokemonWebServer) Start() error {
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

		err = pokeimage.NewPokemonImage(image).Write(w)
		if err != nil {
			log.Printf("error writing to response: %v", err)
		}
	})

	log.Printf("webserver started at %s\n", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
