package internal

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Quadtree struct {
	Bounds   rl.Rectangle
	Objects  []*Object
	Nodes    [4]*Quadtree
	Capacity int
	Divided  bool
}

func (q *Quadtree) Split() {

	halfWidth := q.Bounds.Width / 2
	halfHeight := q.Bounds.Height / 2

	q.Nodes[0] = &Quadtree{
		Bounds: rl.NewRectangle(
			q.Bounds.X,
			q.Bounds.Y,
			halfWidth,
			halfHeight,
		),
		Capacity: q.Capacity,
	}

	q.Nodes[1] = &Quadtree{
		Bounds: rl.NewRectangle(
			q.Bounds.X+halfWidth,
			q.Bounds.Y,
			halfWidth,
			halfHeight,
		),
		Capacity: q.Capacity,
	}

	q.Nodes[2] = &Quadtree{
		Bounds: rl.NewRectangle(
			q.Bounds.X,
			q.Bounds.Y+halfHeight,
			halfWidth,
			halfHeight,
		),
		Capacity: q.Capacity,
	}

	q.Nodes[3] = &Quadtree{
		Bounds: rl.NewRectangle(
			q.Bounds.X+halfWidth,
			q.Bounds.Y+halfHeight,
			halfWidth,
			halfHeight,
		),
		Capacity: q.Capacity,
	}

	q.Divided = true
}

func (q *Quadtree) Insert(object *Object) bool {

	// If leaf, and point outside the segment
	if !rl.CheckCollisionPointRec(object.Pos, q.Bounds) {
		return false
	}

	if len(q.Objects) < q.Capacity {
		q.Objects = append(q.Objects, object)
		return true
	}

	if !q.Divided {
		q.Split()
	}

	for _, node := range q.Nodes {
		if node.Insert(object) {
			return true
		}
	}

	return false

}

func (q *Quadtree) Draw() {
	rl.DrawRectangleLines(
		q.Bounds.ToInt32().X,
		q.Bounds.ToInt32().Y,
		q.Bounds.ToInt32().Width,
		q.Bounds.ToInt32().Height,
		rl.White,
	)

	if q.Divided {
		for _, node := range q.Nodes {
			node.Draw()
		}
	}
}
