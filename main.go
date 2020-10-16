package main

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var (
	red   = color.RGBA{R: 255, A: 255}
	green = color.RGBA{G: 255, A: 255}
	black = color.RGBA{A: 255}
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

const (
	width  = 800
	height = 800
)

func main() {
	model := Open("obj/african_head.obj")

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, black)
		}
	}
	/*for _, face := range model.faces {
		screenCoords := [3]Vec2i{}
		for j := 0; j < 3; j++ {
			worldCoords := model.verts[face[j]]
			screenCoords[j] = Vec2i{int((worldCoords.x + 1) * width / 2), int((worldCoords.y + 1) * height / 2)}
		}
		triangle2(screenCoords[:], img, color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
	}*/
	lightDir := Vec3f{0, 0, -1}
	for _, face := range model.faces {
		screenCoords := [3]Vec2i{}
		worldCoords := [3]Vec3f{}
		for j := 0; j < 3; j++ {
			v := model.verts[face[j]]
			screenCoords[j] = Vec2i{int((v.x + 1) * width / 2.0), int((v.y + 1) * height / 2.0)}
			worldCoords[j] = v
		}
		n := (worldCoords[2].Sub(worldCoords[0])).Cross(worldCoords[1].Sub(worldCoords[0]))
		n = n.Normalize()
		intensity := n.Mul(lightDir).AsFloat()
		if intensity > 0 {
			triangle2(screenCoords[:], img, color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255})
		}
	}
	/*t0 := [3]Vec2i{{10, 70}, {50, 160}, {70, 80}}
	t1 := [3]Vec2i{{180, 50}, {150, 1}, {70, 180}}
	t2 := [3]Vec2i{{180, 150}, {120, 160}, {130, 180}}
	triangle(t0[0], t0[1], t0[2], img, red)
	triangle(t1[0], t1[1], t1[2], img, white)
	triangle(t2[0], t2[1], t2[2], img, green)*/
	//pts := []Vec2i{{10, 10}, {100, 30}, {190, 160}}
	//triangle2(pts, img, red)
	img2 := imaging.FlipV(img)
	f, _ := os.Create("lesson2.png")
	defer f.Close()
	png.Encode(f, img2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func swapV2i(a, b Vec2i) (Vec2i, Vec2i) {
	return b, a
}

//swap returns the two ints but in the swapped order they were sent
func swap(a, b int) (int, int) {
	return b, a
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func barycentric(pts []Vec2i, P Vec2i) Vec3f {
	u := Vec3f{float32(pts[2].x - pts[0].x), float32(pts[1].x - pts[0].x), float32(pts[0].x - P.x)}.Cross(
		Vec3f{float32(pts[2].y - pts[0].y), float32(pts[1].y - pts[0].y), float32(pts[0].y - P.y)})
	/* `pts` and `P` has integer value as coordinates
	so `abs(u[2])` < 1 means `u[2]` is 0, that means
	triangle is degenerate, in this case return something with negative coordinates */
	if math.Abs(float64(u.z)) < 1 {
		return Vec3f{-1, 1, 1}
	}
	return Vec3f{1 - (u.x+u.y)/u.z, u.y / u.z, u.x / u.z}
}

func triangle2(pts []Vec2i, img *image.RGBA, color color.Color) {
	bboxmin := Vec2i{img.Rect.Dx() - 1, img.Rect.Dy() - 1}
	bboxmax := Vec2i{}
	clamp := Vec2i{img.Rect.Dx() - 1, img.Rect.Dy() - 1}
	for i := 0; i < 3; i++ {
		bboxmin.x = max(0, min(bboxmin.x, pts[i].x))
		bboxmax.x = min(clamp.x, max(bboxmax.x, pts[i].x))
		bboxmin.y = max(0, min(bboxmin.y, pts[i].y))
		bboxmax.y = min(clamp.y, max(bboxmax.y, pts[i].y))
	}
	var P Vec2i
	for P.x = bboxmin.x; P.x <= bboxmax.x; P.x++ {
		for P.y = bboxmin.y; P.y <= bboxmax.y; P.y++ {
			bcScreen := barycentric(pts, P)
			if bcScreen.x < 0 || bcScreen.y < 0 || bcScreen.z < 0 {
				continue
			}
			img.Set(P.x, P.y, color)
		}
	}
}

func triangle(t0, t1, t2 Vec2i, img *image.RGBA, color color.Color) {
	if t0.y == t1.y && t0.y == t2.y {
		return // I dont care about degenerate triangles
	}

	// sort the vertices, t0, t1, t2 lower−to−upper (bubblesort yay!)
	if t0.y > t1.y {
		t0, t1 = swapV2i(t0, t1)
	}
	if t0.y > t2.y {
		t0, t2 = swapV2i(t0, t2)
	}
	if t1.y > t2.y {
		t1, t2 = swapV2i(t1, t2)
	}
	totalHeight := float32(t2.y - t0.y)
	for y := t0.y; y <= t1.y; y++ {
		segH := float32(t1.y - t0.y + 1)
		alpha := float32(y-t0.y) / totalHeight
		beta := float32(y-t0.y) / segH
		A := t0.Add(t2.Sub(t0).Scale(alpha))
		B := t0.Add(t1.Sub(t0).Scale(beta))
		if A.x > B.x {
			A, B = swapV2i(A, B)
		}
		for j := A.x; j <= B.x; j++ {
			img.Set(j, y, color) // attention, due to int casts t0.y+i != A.y
		}
	}
	for y := t1.y; y <= t2.y; y++ {
		segH := float32(t2.y - t1.y + 1)
		alpha := float32(y-t0.y) / totalHeight
		beta := float32(y-t1.y) / segH
		A := t0.Add(t2.Sub(t0).Scale(alpha))
		B := t1.Add(t2.Sub(t1).Scale(beta))
		if A.x > B.x {
			A, B = swapV2i(A, B)
		}
		for j := A.x; j <= B.x; j++ {
			img.Set(j, y, color) // attention, due to int casts t0.y+i != A.y
		}
	}
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
