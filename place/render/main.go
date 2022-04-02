package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

func main() {
	conf := *getConfig()
	inputPtr := OpenImage(conf.InputImage)
	if inputPtr == nil {
		panic("did not find input image")
	}
	input := *inputPtr
	inputWidth := input.Bounds().Dx()
	inputHeight := input.Bounds().Dy()
	img := image.NewRGBA(image.Rect(0, 0, 3000, 3000))

	existingPtr := OpenImage(conf.OutputImage)
	if existingPtr != nil {
		existing := *existingPtr
		// copy over
		for x := 0; x < existing.Bounds().Dx(); x++ {
			for y := 0; y < existing.Bounds().Dy(); y++ {
				color := existing.At(x, y)
				img.Set(x, y, color)
			}
		}
	}

	for x := 0; x < inputWidth; x++ {
		for y := 0; y < inputHeight; y++ {
			color := input.At(x, y)
			targetX := (conf.X+x)*3 + 1
			targetY := (conf.Y-inputHeight+1+y)*3 + 1
			img.Set(targetX, targetY, color)
		}
	}
	img.SetRGBA(0, 0, color.RGBA{
		255, 0, 0, 255,
	})
	img.SetRGBA(1, 1, color.RGBA{
		255, 0, 0, 255,
	})
	img.SetRGBA(2, 2, color.RGBA{
		255, 0, 0, 255,
	})
	SaveImage(conf.OutputImage, img)
}

type Config struct {
	InputImage  string
	OutputImage string
	// left first pixel of image
	X int
	// bottom left first pixel of image, y increases down
	Y int
}

func getConfig() *Config {
	return &Config{
		InputImage:  mustReadEnv("INPUT_IMAGE"),
		OutputImage: mustReadEnv("OUTPUT_IMAGE"),
		X:           mustReadIntEnv("X"),
		Y:           mustReadIntEnv("Y"),
	}
}

func mustReadEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("env not found: %s", name))
	}
	return value
}

func mustReadIntEnv(name string) int {
	value := mustReadEnv(name)
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(parsed)
}

func OpenImage(filename string) *image.Image {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(err)
	}
	ret, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return &ret
}

func SaveImage(filename string, img *image.RGBA) {
	f, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
