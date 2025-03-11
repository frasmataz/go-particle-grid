package internal

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Object struct {
	Id     int
	Pos    rl.Vector2
	Vel    rl.Vector2
	Colour color.RGBA
	Mass   float32
}
