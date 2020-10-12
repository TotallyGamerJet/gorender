package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Model struct {
	verts []Vec3f
	faces [][]int
}

func Open(name string) Model {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var line string
	scanner := bufio.NewScanner(f)
	var verts = make([]Vec3f, 0, 25)
	var faces = make([][]int, 0, 25)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "v ") {
			splits := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(splits[1], 32)
			y, _ := strconv.ParseFloat(splits[2], 32)
			z, _ := strconv.ParseFloat(splits[3], 32)
			vec := Vec3f{float32(x), float32(y), float32(z)}
			verts = append(verts, vec)
		} else if strings.HasPrefix(line, "f ") {
			splits := strings.Split(line, " ")
			f := make([]int, 0, 3)
			for _, ff := range splits[1:] {
				fs := strings.Split(ff, "/")
				i, _ := strconv.ParseInt(fs[0], 10, 64)
				i--
				f = append(f, int(i))
			}
			faces = append(faces, f)
		}
	}
	return Model{verts, faces}
}
