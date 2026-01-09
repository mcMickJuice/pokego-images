package pokemon

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
const POKEMON_DETAIL_URL = "https://pokeapi.co/api/v2/pokemon/%s"

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
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))

	if err != nil {
		return pokemonResponse{}, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()
	des := json.NewDecoder(resp.Body)

	var pokemonResp pokemonResponse
	if err := des.Decode(&pokemonResp); err != nil {
		return pokemonResponse{}, err
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
		// wrap this to be specific that image decoding failed
		return nil, fmt.Errorf("failed to decode pokemon sprite: %w", err)
	}

	return img, nil
}
