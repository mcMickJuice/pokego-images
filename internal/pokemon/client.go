package pokemon

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

const (
	pokemonDetailUrl = "https://pokeapi.co/api/v2/pokemon/%s"
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

type PokemonClient struct {
	pokemonName string
}

func NewPokemonClient(pokemonName string) *PokemonClient {

	return &PokemonClient{
		pokemonName: pokemonName,
	}
}

var ErrPokemonNotFound = fmt.Errorf("pokemon not found")

func (pc PokemonClient) GetPokemonSprite() (image.Image, error) {

	pokeResp, err := getPokemon(pc.pokemonName)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon details: %w", err)
	}

	pokeSpriteResp, err := getPokemonSprite(pokeResp.Sprites.Default)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon sprite: %w", err)
	}

	return pokeSpriteResp, nil
}

func getPokemon(pokemonName string) (pokemonResponse, error) {
	resp, err := http.Get(fmt.Sprintf(pokemonDetailUrl, pokemonName))

	if err != nil {
		return pokemonResponse{}, fmt.Errorf("unknown error fetching pokemon: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return pokemonResponse{}, ErrPokemonNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return pokemonResponse{}, fmt.Errorf("failed to fetch pokemon: status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	des := json.NewDecoder(resp.Body)

	var pokemonResp pokemonResponse
	if err := des.Decode(&pokemonResp); err != nil {
		return pokemonResponse{}, fmt.Errorf("failed to decode pokemon response: %w", err)
	}

	return pokemonResp, nil
}

func getPokemonSprite(pokemonSprintUrl string) (image.Image, error) {

	resp, err := http.Get(pokemonSprintUrl)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pokemon sprite: %w", err)
	}

	return img, nil
}
