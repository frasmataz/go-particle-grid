package internal

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Object struct {
	Pos    rl.Vector2
	Vel    rl.Vector2
	Colour color.RGBA
}
