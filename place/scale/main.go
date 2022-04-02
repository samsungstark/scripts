package main

import (
	"fmt"
	"image"
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
	img := image.NewRGBA(image.Rect(0, 0, conf.DesiredWidth, conf.DesiredHeight))
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	for x := 0; x < imgWidth; x++ {
		inputX := int((float64(x) + .5) / float64(imgWidth) * float64(inputWidth))
		for y := 0; y < imgHeight; y++ {
			inputY := int((float64(y) + .5) / float64(imgHeight) * float64(inputHeight))
			color := input.At(inputX, inputY)
			img.Set(x, y, color)
		}
	}
	SaveImage(conf.OutputImage, img)
}

type Config struct {
	InputImage    string
	OutputImage   string
	DesiredWidth  int
	DesiredHeight int
}

func getConfig() *Config {
	return &Config{
		InputImage:    mustReadEnv("INPUT_IMAGE"),
		OutputImage:   mustReadEnv("OUTPUT_IMAGE"),
		DesiredWidth:  mustReadIntEnv("DESIRED_WIDTH"),
		DesiredHeight: mustReadIntEnv("DESIRED_HEIGHT"),
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
