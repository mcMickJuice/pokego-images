package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcmickjuice/pokego/internal/asciiimage"
	"mcmickjuice/pokego/internal/pokemon"
	"net/http"
)

type PokemonWebServer struct {
	addr          string
	pokemonClient *pokemon.PokemonClient
}

func NewPokemonWebServer(addr string) PokemonWebServer {
	return PokemonWebServer{addr, pokemon.NewPokemonClient("https://pokeapi.co")}
}

func (s PokemonWebServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/pokemon/all", func(w http.ResponseWriter, r *http.Request) {
		resp, err := s.pokemonClient.GetPokemonList(r.Context())

		fmt.Println("list has been received")
		if err != nil {

			if errors.Is(err, pokemon.ErrPokemonListNotFound) {

				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, "Pokemon List not Found")
				return
			}

			log.Printf("error getting pokemon list: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Internal server error occurred")
			return
		}

		fmt.Println("prepping responses")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Internal server error occurred")
		}

	})
	mux.HandleFunc("/pokemon/{pokemon}", func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("pokemon")

		resp, err := s.pokemonClient.GetPokemon(r.Context(), param)

		if err != nil {
			if errors.Is(err, pokemon.ErrPokemonNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, "Pokemon %s Not Found", param)
				return
			}
			log.Printf("error getting sprite for %s: %v", param, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "Internal server error occurred")
			return
		}

		image, err := s.pokemonClient.GetPokemonSpriteImage(r.Context(), resp)

		if err != nil {
			if errors.Is(err, pokemon.ErrPokemonSpriteNotFound) {

				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprintf(w, "Pokemon %s sprite not found", param)
				return
			}
			log.Printf("error getting sprite for %s: %v", param, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Internal server error occurred")
			return
		}

		err = asciiimage.NewAsciiImage(image).Write(w)
		if err != nil {
			log.Printf("error writing to response: %v", err)
		}
	})

	log.Printf("webserver started at %s\n", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
