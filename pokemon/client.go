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

type PokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites PokemonSpritesResponse `json:"sprites"`
}

type PokemonSpritesResponse struct {
	Default     string `json:"front_default"`
	BackDefault string `json:"back_default"`
	FrontShiny  string `json:"front_shiny"`
}

type AllPokemonsResponse struct {
	Results []AllPokemonResponse `json:"results"`
}

type AllPokemonResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
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

	pokemonResponse, err := getPokemon(pc.pokemonName)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon details: %w", err)
	}

	pokemonSprite, err := getPokemonSprite(pokemonResponse.Sprites.Default)

	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon sprite: %w", err)
	}

	return pokemonSprite, nil
}

func getPokemon(pokemonName string) (PokemonResponse, error) {
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))

	if err != nil {
		return PokemonResponse{}, err
	}

	defer resp.Body.Close()
	des := json.NewDecoder(resp.Body)

	var pokemonResponse PokemonResponse
	if err := des.Decode(&pokemonResponse); err != nil {
		return PokemonResponse{}, err
	}

	return pokemonResponse, nil
}

func getPokemonSprite(pokemonSprintUrl string) (image.Image, error) {

	resp, err := http.Get(pokemonSprintUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	img, err := png.Decode(resp.Body)
	if err != nil {
		// wrap this to be specific that image decoding failed
		return nil, fmt.Errorf("failed to decode pokemon sprite: %w", err)
	}

	return img, nil
}
