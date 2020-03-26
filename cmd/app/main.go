package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

// Vec3f ...
type Vec3f struct {
	v0 float64
	v1 float64
	v2 float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func clip(num float64) float64 {
	if num <= 0 {
		return 0.
	}

	if num > 1. {
		return 1.
	}

	return num
}

func interpolate(start Vec3f, end Vec3f, percent float64) Vec3f {
	ft := float64(percent * 3.1415927)
	f := (1 - math.Cos(ft)) * 0.5

	return Vec3f{
		clip(start.v0*(1.-f) + end.v0*f),
		clip((start.v1*(1.-f) + end.v1*f)),
		clip((start.v2*(1.-f) + end.v2*f)),
	}
}

func hexRgbToVec(color string) (Vec3f, error) {
	intColor, err := strconv.ParseInt(fmt.Sprintf("ff%s", color), 16, 64)

	if err != nil {
		return Vec3f{0, 0, 0}, err
	}

	r := (intColor & 0x00ff0000) >> 16
	g := (intColor & 0x0000ff00) >> 8
	b := (intColor & 0x0000ff)

	return Vec3f{float64(r / 255), float64(g / 255), float64(b / 255)}, nil
}

var (
	width  int
	height int
	from   string
	to     string
	out    string
)

func init() {
	flag.IntVar(&width, "w", 1000, "output image width")
	flag.IntVar(&height, "h", 1000, "output image height")
	flag.StringVar(&from, "f", "ff0000", "gradient start color in RGB hex format")
	flag.StringVar(&to, "t", "0000ff", "gradient end color in RGB hex format")
	flag.StringVar(&out, "o", "out.ppm", "output filename")
}

func main() {
	flag.Parse()

	start, err := hexRgbToVec(from)
	check(err)

	end, err := hexRgbToVec(to)
	check(err)

	framebuffer := make([]Vec3f, width*height)

	for i := 0; i < width*height; i++ {
		framebuffer[i] = interpolate(start, end, float64(i%width)/float64(width))
	}

	f, err := os.Create(out)
	check(err)

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(fmt.Sprintf("P6\n%d %d \n255\n", width, height))
	check(err)

	for i := 0; i < height*width; i++ {
		_ = w.WriteByte(uint8(255 * math.Max(0., math.Min(1., framebuffer[i].v0))))
		_ = w.WriteByte(uint8(255 * math.Max(0., math.Min(1., framebuffer[i].v1))))
		_ = w.WriteByte(uint8(255 * math.Max(0., math.Min(1., framebuffer[i].v2))))
	}

	w.Flush()
}
