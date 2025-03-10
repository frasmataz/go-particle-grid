package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/frasmataz/p5"
)

type ScreenSettings struct {
	height int
	width  int
}

type Config struct {
	initObjectCount int
	screen          ScreenSettings
}

type Object struct {
	x      int
	y      int
	colour color.NRGBA
}

type State struct {
	objects []Object
	config  Config
}

var state State

func main() {
	state = State{
		config: Config{
			initObjectCount: 500000,
			screen: ScreenSettings{
				1024,
				1024,
			},
		},
	}
	p5.Run(setup, draw)
}

func setup() {

	p5.Canvas(state.config.screen.width, state.config.screen.height)

	for _ = range state.config.initObjectCount {
		state.objects = append(state.objects, Object{
			x: rand.Intn(state.config.screen.width),
			y: rand.Intn(state.config.screen.height),
			colour: color.NRGBA{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				255,
			},
		})
	}

}

func draw() {

	image := image.NewNRGBA(
		image.Rectangle{
			image.Point{
				0,
				0,
			},
			image.Point{
				state.config.screen.width,
				state.config.screen.height,
			},
		},
	)

	for _, object := range state.objects {

		image.SetNRGBA(
			object.x,
			object.y,
			object.colour,
		)

	}

	p5.DrawImage(
		image,
		0, 0,
	)
}
