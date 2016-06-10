package lissajous

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"time"
)

var palette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xFF},
	color.RGBA{0x00, 0x00, 0xFF, 0xFF},
	color.RGBA{0x00, 0xFF, 0x00, 0xFF},
	color.RGBA{0x00, 0xFF, 0xFF, 0xFF},
	color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	color.RGBA{0xFF, 0x00, 0xFF, 0xFF},
	color.RGBA{0xFF, 0xFF, 0x00, 0xFF},
	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
}

const backgroundIndex = 0

const (
	Cycles   = 4     // number of complete x oscillator revolutions
	Res      = 0.001 // angular resolution
	Side     = 400   // image canvas side in pixels [0..side]
	NFrames  = 64    // number of animation frames
	Delay    = 8     // delay between frames in 10ms units
	PhaseInc = 0.1   // how much phase to increment in each frame
	FreqDiff = 2.3   // frequency difference between x and y
)

type Conf struct {
	Cycles   int
	Res      float64
	Side     int
	NFrames  int
	Delay    int
	PhaseInc float64
	FreqDiff float64
}

func DefaultConf() *Conf {
	return &Conf{
		Cycles:   Cycles,
		Res:      Res,
		Side:     Side,
		NFrames:  NFrames,
		Delay:    Delay,
		PhaseInc: PhaseInc,
		FreqDiff: FreqDiff,
	}
}

func Gif(out io.Writer, conf *Conf) error {
	var phase float64
	anim := gif.GIF{LoopCount: conf.NFrames}

	for i := 0; i < conf.NFrames; i++ {
		frame, delay := createFrame(anim, conf, phase)
		anim.Image = append(anim.Image, frame)
		anim.Delay = append(anim.Delay, delay)
		phase += conf.PhaseInc
	}

	return gif.EncodeAll(out, &anim)
}

func randomSeedUsingTime() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func createFrame(anim gif.GIF, conf *Conf, phase float64) (*image.Paletted, int) {
	rect := image.Rect(0, 0, conf.Side, conf.Side)
	img := image.NewPaletted(rect, palette)

	for t := 0.0; t < float64(conf.Cycles)*2*math.Pi; t += conf.Res {
		x := math.Sin(t)
		y := math.Sin(t*conf.FreqDiff + phase)

		px, py := cartesianToImage(x, y, conf.Side)
		colorIndex := colorIndexFromT(t, conf.Cycles)
		img.SetColorIndex(px, py, colorIndex)
	}

	return img, conf.Delay
}

func cartesianToImage(x, y float64, side int) (int, int) {
	cX := (x + 1.0) * float64(side) / 2
	cY := (-y + 1.0) * float64(side) / 2

	return int(cX), int(cY)
}

func colorIndexFromT(t float64, cycles int) uint8 {
	index := float64(len(palette)) * t / (float64(cycles) * 2 * math.Pi)
	return uint8(math.Floor(index))
}
