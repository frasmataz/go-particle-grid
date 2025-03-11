package main

import (
	"image/color"
	"math/rand"

	"github.com/frasmataz/go-particle-grid/internal"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScreenSettings struct {
	rect   rl.Rectangle
	width  int32
	height int32
	fps    int32
}

type Config struct {
	initObjectCount int
	screen          ScreenSettings
}

type State struct {
	config      Config
	objects     []internal.Object
	framebuffer rl.Texture2D
}

var state State

func main() {
	state = State{
		config: Config{
			initObjectCount: 500,
			screen: ScreenSettings{
				rect: rl.NewRectangle(
					0, 0, 1000, 800,
				),
				fps: 60,
			},
		},
	}

	rl.InitWindow(
		state.config.screen.rect.ToInt32().Width,
		state.config.screen.rect.ToInt32().Height,
		"particles",
	)
	rl.SetTargetFPS(state.config.screen.fps)

	state.framebuffer = rl.LoadTextureFromImage(
		rl.GenImageColor(
			int(state.config.screen.rect.ToInt32().Width),
			int(state.config.screen.rect.ToInt32().Height),
			rl.Black,
		),
	)

	for range state.config.initObjectCount {
		state.objects = append(state.objects, internal.Object{
			Pos: rl.NewVector2(
				rand.Float32()*state.config.screen.rect.Width,
				rand.Float32()*state.config.screen.rect.Height,
			),
			Vel: rl.NewVector2(
				(rand.Float32()-0.5)*50.0,
				(rand.Float32()-0.5)*50.0,
			),
			Colour: color.RGBA{
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

	update()

	rl.ClearBackground(rl.Black)

	qt := internal.Quadtree{
		Bounds:   state.config.screen.rect,
		Capacity: 4,
	}

	for _, object := range state.objects {
		qt.Insert(&object)
	}

	sampleRect := rl.NewRectangle(200, 200, 200, 200)
	collidedPoints := qt.Query(sampleRect)

	size := state.config.screen.rect.ToInt32().Width * state.config.screen.rect.ToInt32().Height
	pixels := make([]color.RGBA, size)

	for _, object := range state.objects {
		index := (int(object.Pos.Y) * int(state.config.screen.rect.Width)) + int(object.Pos.X)

		if index >= 0 && index < int(size) {
			pixels[index] = color.RGBA(object.Colour)
		}

	}

	for _, object := range collidedPoints {
		index := (int(object.Pos.Y) * int(state.config.screen.rect.Width)) + int(object.Pos.X)

		if index >= 0 && index < int(size) {
			pixels[index] = rl.Green
		}

	}

	rl.BeginDrawing()

	rl.UpdateTexture(state.framebuffer, pixels)
	rl.DrawTexture(state.framebuffer, 0, 0, rl.White)
	qt.Draw()
	rl.DrawRectangleLinesEx(sampleRect, 2.0, rl.Green)

	rl.EndDrawing()

}

func update() {

	for i, object := range state.objects {

		nextPos := rl.Vector2Add(object.Pos, object.Vel)

		//Collide with walls
		if nextPos.X < 0 {
			state.objects[i].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(1, 0))
		} else if nextPos.X >= state.config.screen.rect.Width {
			state.objects[i].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(-1, 0))
		}

		if nextPos.Y < 0 {
			state.objects[i].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(0, 1))
		} else if nextPos.Y >= state.config.screen.rect.Height {
			state.objects[i].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(0, -1))
		}

		state.objects[i].Pos = nextPos

		// Gravity
		state.objects[i].Vel = rl.Vector2Add(state.objects[i].Vel, rl.NewVector2(0, 0.1))

		//Drag
		state.objects[i].Vel = rl.Vector2Scale(state.objects[i].Vel, 0.999)

	}

}
