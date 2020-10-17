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
	blue  = color.RGBA{B: 255, A: 255}
	black = color.RGBA{A: 255}
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

const (
	width  = 800
	height = 800
)

func main() {
	model := OpenOBJ("obj/african_head")

	var zbuffer [width * height]float32
	for i := range zbuffer {
		zbuffer[i] = -math.MaxFloat32
	}
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, black)
		}
	}
	lightDir := Vec3f{0, 0, -1}
	for _, face := range model.faces {
		worldCoords := [3]Vec3f{}
		var pts [3]Vec3f
		for i := range pts {
			v := model.verts[face[i]]
			pts[i] = world2Screen(v)
			worldCoords[i] = v
		}
		n := (worldCoords[2].Sub(worldCoords[0])).Cross(worldCoords[1].Sub(worldCoords[0]))
		n = n.Normalize()
		intensity := n.Mul(lightDir).AsFloat()
		if intensity > 0 {
			triangle2(pts[:], zbuffer[:], img, color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255})
		}
	}
	img2 := imaging.FlipV(img)
	f, _ := os.Create("lesson3.png")
	defer f.Close()
	png.Encode(f, img2)
}

func world2Screen(v Vec3f) Vec3f {
	return Vec3f{float32(int((v.x+1.)*width/2. + .5)), float32(int((v.y+1.)*height/2. + .5)), v.z}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func absf(x float32) float32 {
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

func max(a, b float32) float32 {
	if b > a {
		return b
	}
	return a
}

func min(a, b float32) float32 {
	if b < a {
		return b
	}
	return a
}

func rasterize(p0, p1 Vec2i, img *image.RGBA, color color.Color, ybuffer []int) {
	if p0.x > p1.x {
		p0, p1 = swapV2i(p0, p1)
	}
	for x := p0.x; x <= p1.x; x++ {
		t := float32(x-p0.x) / float32(p1.x-p0.x)
		y := int(float32(p0.y)*(1-t) + float32(p1.y)*t)
		if ybuffer[x] < y {
			ybuffer[x] = y
			img.Set(x, 0, color)
		}
	}
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

func barycentric2(A, B, C, P Vec3f) Vec3f {
	var s [2]Vec3f
	for i := 2; i > 0; {
		i--
		s[i].x = C.Get(i) - A.Get(i)
		s[i].y = B.Get(i) - A.Get(i)
		s[i].z = A.Get(i) - P.Get(i)
	}

	u := s[0].Cross(s[1])
	if absf(u.z) > 1e-2 { // dont forget that u[2] is integer. If it is zero then triangle ABC is degenerate
		return Vec3f{1. - (u.x+u.y)/u.z, u.y / u.z, u.x / u.z}
	}
	return Vec3f{-1, 1, 1} // in this case generate negative coordinates, it will be thrown away by the rasterizer
}

func triangle2(pts []Vec3f, zbuffer []float32, img *image.RGBA, color color.Color) {
	bboxmin := Vec2f{math.MaxFloat32, math.MaxFloat32}
	bboxmax := Vec2f{-math.MaxFloat32, -math.MaxFloat32}
	clamp := Vec2f{float32(img.Rect.Dx() - 1), float32(img.Rect.Dy() - 1)}
	for i := 0; i < 3; i++ {
		bboxmin.x = max(0, min(bboxmin.x, pts[i].x))
		bboxmax.x = min(clamp.x, max(bboxmax.x, pts[i].x))
		bboxmin.y = max(0, min(bboxmin.y, pts[i].y))
		bboxmax.y = min(clamp.y, max(bboxmax.y, pts[i].y))
	}
	var P Vec3f
	for P.x = bboxmin.x; P.x <= bboxmax.x; P.x++ {
		for P.y = bboxmin.y; P.y <= bboxmax.y; P.y++ {
			bcScreen := barycentric2(pts[0], pts[1], pts[2], P)
			if bcScreen.x < 0 || bcScreen.y < 0 || bcScreen.z < 0 {
				continue
			}
			P.z = 0
			for i := 0; i < 3; i++ {
				P.z += pts[i].z * bcScreen.Get(i)
			}
			if zbuffer[int(P.x+P.y*width)] < P.z {
				zbuffer[int(P.x+P.y*width)] = P.z
				img.Set(int(P.x), int(P.y), color)
			}
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

func line2(p0, p1 Vec2i, img *image.RGBA, color color.Color) {
	steep := false
	if abs(p0.x-p1.x) < abs(p0.y-p1.y) {
		p0.x, p0.y = swap(p0.x, p0.y)
		p1.x, p1.y = swap(p1.x, p1.y)
		steep = true
	}
	if p0.x > p1.x {
		p0, p1 = swapV2i(p0, p1)
	}
	for x := p0.x; x <= p1.x; x++ {
		t := float32(x-p0.x) / float32(p1.x-p0.x)
		y := int(float32(p0.y)*(1-t) + float32(p1.y)*t + 0.5)
		if steep {
			img.Set(y, x, color)
		} else {
			img.Set(x, y, color)
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
