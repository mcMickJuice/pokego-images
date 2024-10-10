package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"net/http"
)

type PokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites PokemonSpritesResponse `json:"sprites"`
}

type PokemonSpritesResponse struct {
	Default string `json:"front_default"`
}

type AllPokemonsResponse struct {
	Results []AllPokemonResponse `json:"results"`
}

type AllPokemonResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
const POKEMON_DETAIL_URL = "https://pokeapi.co/api/v2/pokemon/%s"

func main() {
	pokemonName := "snorlax"
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))
	panicIfErr(err)

	defer resp.Body.Close()
	var pokemonResponse PokemonResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&pokemonResponse)

	resp, err = http.Get(pokemonResponse.Sprites.Default)
	panicIfErr(err)
	defer resp.Body.Close()

	// read contents of image response
	image, imageErr := png.Decode(resp.Body)
	panicIfErr(imageErr)

	fmt.Println(image.At(45, 43))
	maxBounds := image.Bounds().Max.X
	var slc = make([]string, maxBounds, maxBounds)

	for r := 0; r < maxBounds; r++ {
		for c := 0; c < maxBounds; c++ {
			grayscale := toGrayscale(image.At(c, r))
			ascii := grayscaleToAscii(grayscale)
			slc[r] += ascii
		}
	}
	for _, line := range slc {
		fmt.Println(line)
	}
}

const R_FACTOR float32 = 0.299
const G_FACTOR float32 = 0.587
const B_FACTOR float32 = 0.114

func toGrayscale(color color.Color) float32 {
	r, g, b, _ := color.RGBA()
	return float32(r)*R_FACTOR + float32(g)*G_FACTOR + float32(b)*B_FACTOR
}

const ASCII = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
const MAX = 65535

// 256*256 - 1 - 0 - 65535
func grayscaleToAscii(brightness float32) string {
	asciiLen := len(ASCII) - 1
	unit := MAX / asciiLen
	inverted := MAX - brightness
	// my god...
	index := inverted / float32(unit)
	indexFloor := math.Floor(float64(index)) // no index out of range!
	return string(ASCII[int64(indexFloor)])
}
