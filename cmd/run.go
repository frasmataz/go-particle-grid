package main

import (
	"image/color"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScreenSettings struct {
	width  int32
	height int32
	fps    int32
}

type Config struct {
	initObjectCount int
	screen          ScreenSettings
}

type Object struct {
	x      int
	y      int
	colour color.RGBA
}

type State struct {
	config      Config
	objects     []Object
	framebuffer rl.Texture2D
}

var state State

func main() {
	state = State{
		config: Config{
			initObjectCount: 500000,
			screen: ScreenSettings{
				width:  1000,
				height: 1000,
				fps:    60,
			},
		},
	}

	rl.InitWindow(
		state.config.screen.width,
		state.config.screen.height,
		"particles",
	)
	rl.SetTargetFPS(state.config.screen.fps)

	state.framebuffer = rl.LoadTextureFromImage(
		rl.GenImageColor(
			int(state.config.screen.width),
			int(state.config.screen.height),
			rl.Black,
		),
	)

	for range state.config.initObjectCount {
		state.objects = append(state.objects, Object{
			x: rand.Intn(int(state.config.screen.width)),
			y: rand.Intn(int(state.config.screen.height)),
			colour: color.RGBA{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				255,
			},
		})
	}

	for !rl.WindowShouldClose() {
		draw()
	}
}

func draw() {

	pixels := make([]color.RGBA, state.config.screen.width*state.config.screen.height)

	for _, object := range state.objects {

		pixels[object.y*int(state.config.screen.width)+object.x] = color.RGBA(object.colour)

	}

	rl.BeginDrawing()

	rl.UpdateTexture(state.framebuffer, pixels)
	rl.DrawTexture(state.framebuffer, 0, 0, rl.White)

	rl.EndDrawing()

}
