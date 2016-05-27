package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var palette = []color.Color{
	color.Black,
	color.RGBA{0x76, 0xEE, 0x00, 0xFF}, // osciloscope green
}

const (
	backgroundIndex = iota
	foregroundIndex
)

const (
	output  = "/tmp/output.gif"
	cycles  = 4     // number of complete x oscillator revolutions
	res     = 0.001 // angular resolution
	side    = 400   // image canvas side in pixels [0..side]
	nframes = 64    // number of animation frames
	delay   = 8     // delay between frames in 10ms units
)

func main() {
	file, err := os.Create(output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	randomSeedUsingTime()

	if err := lissajous(file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("output is at", output)
}

func randomSeedUsingTime() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func lissajous(out io.Writer) error {
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	phase := 0.0                 // phase difference
	anim := gif.GIF{LoopCount: nframes}

	for i := 0; i < nframes; i++ {
		frame, delay := createFrame(anim, freq, phase)
		anim.Image = append(anim.Image, frame)
		anim.Delay = append(anim.Delay, delay)
		phase += 0.1
	}

	return gif.EncodeAll(out, &anim)
}

func createFrame(anim gif.GIF, freq, phase float64) (*image.Paletted, int) {
	rect := image.Rect(0, 0, side, side)
	img := image.NewPaletted(rect, palette)

	for t := 0.0; t < cycles*2*math.Pi; t += res {
		x := math.Sin(t)
		y := math.Sin(t*freq + phase)
		px, py := cartesianToImage(x, y)
		img.SetColorIndex(px, py, foregroundIndex)
	}

	return img, delay
}

func cartesianToImage(x, y float64) (int, int) {
	cX := (x + 1.0) * side / 2
	cY := (-y + 1.0) * side / 2

	return int(cX), int(cY)
}
