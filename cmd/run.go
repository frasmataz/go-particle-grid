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
			initObjectCount: 10000,
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

	for i := range state.config.initObjectCount {
		state.objects = append(state.objects, internal.Object{
			Id: i,
			Pos: rl.NewVector2(
				rand.Float32()*state.config.screen.rect.Width,
				rand.Float32()*state.config.screen.rect.Height,
			),
			Vel: rl.NewVector2(
				(rand.Float32()-0.5)*5.0,
				(rand.Float32()-0.5)*10.0,
			),
			Colour: color.RGBA{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				255,
			},
			Mass: 1.0,
		})
	}

	for !rl.WindowShouldClose() {
		draw()
	}
}

func draw() {

	update()

	rl.ClearBackground(rl.Black)

	size := state.config.screen.rect.ToInt32().Width * state.config.screen.rect.ToInt32().Height
	pixels := make([]color.RGBA, size)

	for _, object := range state.objects {
		index := (int(object.Pos.Y) * int(state.config.screen.rect.Width)) + int(object.Pos.X)

		if index >= 0 && index < int(size) {
			pixels[index] = color.RGBA(object.Colour)
		}
	}

	rl.BeginDrawing()

	rl.UpdateTexture(state.framebuffer, pixels)
	rl.DrawTexture(state.framebuffer, 0, 0, rl.White)

	rl.EndDrawing()

}

func update() {

	// Build quadtree
	qt := internal.Quadtree{
		Bounds:   state.config.screen.rect,
		Capacity: 4,
	}

	for _, object := range state.objects {
		qt.Insert(&object)
	}

	// Update each object
	for oi, object := range state.objects {

		nextPos := rl.Vector2Add(object.Pos, object.Vel)

		nextRect := rl.NewRectangle(
			float32(int(nextPos.X)),
			float32(int(nextPos.Y)),
			1.0,
			1.0,
		)

		points := qt.Query(nextRect)

		for _, point := range points {
			if point.Id != object.Id {

				// Elastic collisions - ported from https://www.plasmaphysics.org.uk/programs/coll2d_cpp.htm
				// Precalc some values
				m21 := point.Mass / object.Mass
				x21 := point.Pos.X - object.Pos.X
				y21 := point.Pos.Y - object.Pos.Y
				vx21 := point.Vel.X - object.Vel.X
				vy21 := point.Vel.Y - object.Vel.Y

				// If balls are approaching:
				if (vx21*x21 + vy21*y21) < 0 {
					// Calculate bounce
					a := y21 / x21
					dvx2 := -2 * (vx21 + a*vy21) / ((1 + a*a) * (1 + m21))
					vx2 := point.Vel.X + dvx2
					vy2 := point.Vel.Y + a*dvx2
					vx1 := object.Vel.X - m21*dvx2
					vy1 := object.Vel.Y - a*m21*dvx2

					// Set resulting new velocities
					state.objects[oi].Vel = rl.NewVector2(vx1, vy1)
					point.Vel = rl.NewVector2(vx2, vy2)
					//
					// // Work out proportion of damping according to relative masses
					// ball1DampingVec := vek.DivNumber(vek.MulNumber(ball.velocity, mass1), (mass1 + mass2))
					// ball2DampingVec := vek.DivNumber(vek.MulNumber(ball2.velocity, mass2), (mass1 + mass2))
					//
					// // Apply damping vectors scaled by bounciness parameter
					// ball.velocity =
					// 	vek.Add(
					// 		vek.MulNumber(
					// 			vek.Sub(ball.velocity, ball1DampingVec),
					// 			bounciness),
					// 		ball1DampingVec)
					// ball2.velocity = vek.Add(vek.MulNumber(vek.Sub(ball2.velocity, ball2DampingVec), bounciness), ball2DampingVec)
				}

			}
		}

		//Collide with walls
		if nextPos.X < 0 {
			state.objects[oi].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(1, 0))
			nextPos.X += 1.0
		} else if nextPos.X >= state.config.screen.rect.Width-1 {
			state.objects[oi].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(-1, 0))
			nextPos.X -= 1.0
		}

		if nextPos.Y < 0 {
			state.objects[oi].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(0, 1))
			nextPos.Y += 1.0
		} else if nextPos.Y >= state.config.screen.rect.Height-1 {
			state.objects[oi].Vel = rl.Vector2Reflect(object.Vel, rl.NewVector2(0, -1))
			nextPos.Y -= 1.0
		}

		state.objects[oi].Pos = nextPos

		// Gravity
		state.objects[oi].Vel = rl.Vector2Add(state.objects[oi].Vel, rl.NewVector2(0, 0.1))

		//Drag
		state.objects[oi].Vel = rl.Vector2Scale(state.objects[oi].Vel, 0.999)

	}

}
