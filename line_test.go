package main

import (
	"image"
	"image/color"
	"testing"
)

func Benchmark_Line(b *testing.B) {
	const width = 100
	const height = 100

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	white := color.RGBA{}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		line(13, 20, 80, 40, img, white)
	}
}
