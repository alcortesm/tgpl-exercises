package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"os"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = iota
	blackIndex
)

const (
	output  = "/tmp/output.gif"
	cycles  = 1     // number of complete x oscillator revolutions
	res     = 0.001 // angular resolution
	side    = 200   // image canvas side in pixels [0..side]
	nframes = 1     // number of animation frames
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

	if err := lissajous(file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("output is at", output)
}

func lissajous(out io.Writer) error {
	//freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	freq := 1.0  // relative frequency of y oscillator
	phase := 0.0 // phase difference
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
		cX, cY := cartesianToImage(x, y)
		img.SetColorIndex(cX, cY, blackIndex)
	}

	return img, delay
}

func cartesianToImage(x, y float64) (int, int) {
	cX := (x + 1.0) * side / 2
	cY := (-y + 1.0) * side / 2

	return int(cX), int(cY)
}
