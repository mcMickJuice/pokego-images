package pokeimage

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type PokemonImage struct {
	image.Image
}

func NewPokemonImage(img image.Image) *PokemonImage {
	return &PokemonImage{
		Image: img,
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

func (pi PokemonImage) ToAsciiArt() {

	maxBounds := pi.Bounds().Max.X
	var slc = make([]string, maxBounds)

	for r := 0; r < maxBounds; r++ {
		for c := 0; c < maxBounds; c++ {
			grayscale := toGrayscale(pi.At(c, r))
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
