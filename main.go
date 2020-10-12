package main

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	model := Open("obj/african_head.obj")
	const width = 800
	const height = 800

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	//red := color.RGBA{R: 255, A: 255}
	black := color.RGBA{A: 255}
	white := color.RGBA{}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, black)
		}
	}
	for _, face := range model.faces {
		for j := 0; j < 3; j++ {
			v0 := model.verts[face[j]]
			v1 := model.verts[face[(j+1)%3]]
			x0 := int((v0.x + 1) * width / 2.0)
			y0 := int((v0.y + 1.) * height / 2.0)
			x1 := int((v1.x + 1.) * width / 2.0)
			y1 := int((v1.y + 1.) * height / 2.0)
			line(x0, y0, x1, y1, img, white)
		}
	}
	img2 := imaging.FlipV(img)
	f, _ := os.Create("lesson1.png")
	defer f.Close()
	png.Encode(f, img2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//swap returns the two ints but in the swapped order they were sent
func swap(a, b int) (int, int) {
	return b, a
}

func line(x0, y0, x1, y1 int, img *image.RGBA, color color.Color) {
	steep := false
	if abs(x0-x1) < abs(y0-y1) {
		x0, y0 = swap(x0, y0)
		x1, y1 = swap(x1, y1)
		steep = true
	}
	if x0 > x1 {
		x0, x1 = swap(x0, x1)
		y0, y1 = swap(y0, y1)
	}
	if steep {
		for x := x0; x <= x1; x++ {
			t := float64(x-x0) / float64(x1-x0)
			y := int(float64(y0)*(1-t) + float64(y1)*t)
			img.Set(y, x, color)
		}
	} else {
		for x := x0; x <= x1; x++ {
			t := float64(x-x0) / float64(x1-x0)
			y := int(float64(y0)*(1-t) + float64(y1)*t)
			img.Set(x, y, color)
		}
	}
}
