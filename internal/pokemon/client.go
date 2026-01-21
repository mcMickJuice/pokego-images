package pokemon

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

const (
	pokemonDetailURL = "/api/v2/pokemon/%s"
)

type pokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites pokemonSpritesResponse `json:"sprites"`
}

type pokemonSpritesResponse struct {
	Default     string `json:"front_default"`
	BackDefault string `json:"back_default"`
	FrontShiny  string `json:"front_shiny"`
}

// Use to fetch Pokémon data from the PokeAPI.
type PokemonClient struct {
	pokeApiBaseURL string
}

func NewPokemonClient(pokeApiBaseURL string) *PokemonClient {
	return &PokemonClient{pokeApiBaseURL}
}

var ErrPokemonNotFound = fmt.Errorf("pokemon not found")

func (pc PokemonClient) GetPokemonSprite(pokemonName string) (image.Image, error) {

	pokeResp, err := pc.getPokemon(pokemonName)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon details: %w", err)
	}

	pokeSpriteResp, err := getPokemonSprite(pokeResp.Sprites.Default)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon sprite: %w", err)
	}

	return pokeSpriteResp, nil
}

func (pc PokemonClient) getPokemon(pokemonName string) (pokemonResponse, error) {
	url := pc.pokeApiBaseURL + pokemonDetailURL
	resp, err := http.Get(fmt.Sprintf(url, pokemonName))

	if err != nil {
		return pokemonResponse{}, fmt.Errorf("unknown error fetching pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return pokemonResponse{}, ErrPokemonNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return pokemonResponse{}, fmt.Errorf("failed to fetch pokemon: status code %v", resp.StatusCode)
	}

	des := json.NewDecoder(resp.Body)

	var pokemonResp pokemonResponse
	if err := des.Decode(&pokemonResp); err != nil {
		return pokemonResponse{}, fmt.Errorf("failed to decode pokemon response: %w", err)
	}

	return pokemonResp, nil
}

func getPokemonSprite(pokemonSpriteURL string) (image.Image, error) {

	resp, err := http.Get(pokemonSpriteURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pokemon sprite: %w", err)
	}

	return img, nil
}
