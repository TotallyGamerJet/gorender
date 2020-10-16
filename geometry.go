package main

import "math"

type Vec3f struct {
	x, y, z float32
}

type Vec2i struct {
	x, y int
}

func (v Vec2i) Add(v2 Vec2i) Vec2i {
	return Vec2i{v.x + v2.x, v.y + v2.y}
}

func (v Vec2i) Sub(v2 Vec2i) Vec2i {
	return Vec2i{v.x - v2.x, v.y - v2.y}
}

func (v Vec2i) Scale(i float32) Vec2i {
	return Vec2i{int(float32(v.x) * i), int(float32(v.y) * i)}
}

func (v Vec3f) Cross(v2 Vec3f) Vec3f {
	return Vec3f{v.y*v2.z - v.z*v2.y, v.z*v2.x - v.x*v2.z, v.x*v2.y - v.y*v2.x}
}

func (v Vec3f) Add(v2 Vec3f) Vec3f {
	return Vec3f{v.x + v2.x, v.y + v2.y, v.z + v2.z}
}

func (v Vec3f) Sub(v2 Vec3f) Vec3f {
	return Vec3f{v.x - v2.x, v.y - v2.y, v.z - v2.z}
}

func (v Vec3f) Mul(v2 Vec3f) Vec3f {
	return Vec3f{v.x * v2.x, v.y * v2.y, v.z * v2.z}
}

func (v Vec3f) Scale(i float32) Vec3f {
	return Vec3f{v.x * i, v.y * i, v.z * i}
}

func (v Vec3f) AsFloat() float32 {
	return v.x + v.y + v.z
}

//float norm () const { return std::sqrt(x*x+y*y+z*z); }
func (v Vec3f) norm() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

func (v Vec3f) Normalize() Vec3f {
	return v.Scale(1 / v.norm())
}
