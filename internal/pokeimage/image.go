package pokeimage

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
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

const (
	rFactor            float32 = 0.299
	gFactor            float32 = 0.587
	bFactor            float32 = 0.114
	asciiChars         string  = " `.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@"
	maxAsciiRange              = 65535
	asciiSpaceCharCode         = 32
)

func toGrayscale(color color.Color) float32 {
	r, g, b, _ := color.RGBA()
	return float32(r)*rFactor + float32(g)*gFactor + float32(b)*bFactor
}

// a better implementation is to use bytes.TrimSpace(line) then check for length.
// make this change after we have unit tests in place
func blankLine(line []byte) bool {
	isBlankLine := true
	for _, char := range line {
		if char != asciiSpaceCharCode {
			isBlankLine = false
			break
		}
	}
	return isBlankLine
}

// 256*256 - 1 - 0 - 65535
func grayscaleToAscii(brightness float32) byte {
	asciiLen := len(asciiChars) - 1
	unit := maxAsciiRange / asciiLen
	// my god...
	index := brightness / float32(unit)
	indexFloor := math.Floor(float64(index)) // no index out of range!
	return (asciiChars[int64(indexFloor)])
}

func (pi PokemonImage) Write(w io.Writer) error {
	maxBounds := pi.Bounds().Max.X
	var slc = make([][]byte, maxBounds)

	// this assumes images are square. for row count, use pi.Bounds().Max.Y
	// make this change after we have unit tests in place
	for r := 0; r < maxBounds; r++ {
		var buf bytes.Buffer
		for c := 0; c < maxBounds; c++ {
			grayscale := toGrayscale(pi.At(c, r))
			ascii := grayscaleToAscii(grayscale)
			buf.WriteByte(ascii)
		}
		slc[r] = buf.Bytes()
	}

	for _, line := range slc {
		if !blankLine(line) {
			// append newline here instead of above so as not to break blaneLine logic
			line = append(line, '\n')
			_, err := w.Write(line)
			if err != nil {
				return fmt.Errorf("error writing to writer: %w", err)
			}
		}
	}
	return nil

}
