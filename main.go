package main

import (
	"encoding/json"
	"flag"
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
	pokemonPtr := flag.String("pokemon", "snorlax", "give a pokemon")
	flag.Parse()
	resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, *pokemonPtr))
	panicIfErr(err)

	defer resp.Body.Close()
	var pokemonResponse PokemonResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&pokemonResponse); err != nil {
		panicIfErr(err)
	}

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
		if !blankLine(line) {

			fmt.Println(line)
		}
	}
}

const R_FACTOR float32 = 0.299
const G_FACTOR float32 = 0.587
const B_FACTOR float32 = 0.114

func toGrayscale(color color.Color) float32 {
	r, g, b, _ := color.RGBA()
	return float32(r)*R_FACTOR + float32(g)*G_FACTOR + float32(b)*B_FACTOR
}

func blankLine(line string) bool {
	isBlankLine := true
	for _, char := range line {
		if char != 32 {
			isBlankLine = false
			break
		}
	}
	return isBlankLine
}

const ASCII = " `.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@"

const MAX = 65535

// 256*256 - 1 - 0 - 65535
func grayscaleToAscii(brightness float32) string {
	asciiLen := len(ASCII) - 1
	unit := MAX / asciiLen
	// my god...
	index := brightness / float32(unit)
	indexFloor := math.Floor(float64(index)) // no index out of range!
	return string(ASCII[int64(indexFloor)])
}
