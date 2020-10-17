package main

import (
	"bufio"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
)

type Model struct {
	verts   []Vec3f
	faces   [][]int
	norms   []Vec3f
	uvs     []Vec2f
	diffuse image.Image
}

func OpenOBJ(name string) Model {
	f, err := os.Open(name + ".obj")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var line string
	scanner := bufio.NewScanner(f)
	var verts = make([]Vec3f, 0, 25)
	var faces = make([][]int, 0, 25)
	var norms = make([]Vec3f, 0, 25)
	var uvs = make([]Vec2f, 0, 25)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "v ") {
			splits := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(splits[1], 32)
			y, _ := strconv.ParseFloat(splits[2], 32)
			z, _ := strconv.ParseFloat(splits[3], 32)
			vec := Vec3f{float32(x), float32(y), float32(z)}
			verts = append(verts, vec)
		} else if strings.HasPrefix(line, "vn ") {
			splits := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(splits[1], 32)
			y, _ := strconv.ParseFloat(splits[2], 32)
			z, _ := strconv.ParseFloat(splits[3], 32)
			vec := Vec3f{float32(x), float32(y), float32(z)}
			norms = append(norms, vec)
		} else if strings.HasPrefix(line, "vt ") {
			splits := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(splits[1], 32)
			y, _ := strconv.ParseFloat(splits[2], 32)
			uv := Vec2f{float32(x), float32(y)}
			uvs = append(uvs, uv)
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
	diffuse := loadTexture(name + "_diffuse.png")
	return Model{verts, faces, norms, uvs, diffuse}
}

func loadTexture(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}

//void Model::load_texture(std::string filename, const char *suffix, TGAImage &img) {
//    std::string texfile(filename);
//    size_t dot = texfile.find_last_of(".");
//    if (dot!=std::string::npos) {
//        texfile = texfile.substr(0,dot) + std::string(suffix);
//        std::cerr << "texture file " << texfile << " loading " << (img.read_tga_file(texfile.c_str()) ? "ok" : "failed") << std::endl;
//        img.flip_vertically();
//    }
//}
//
//TGAColor Model::diffuse(Vec2i uv) {
//    return diffusemap_.get(uv.x, uv.y);
//}
