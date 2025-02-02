/*
Package fov implements basic recursive shadowcasting for displaying a field of view on a 2D Grid
The exact structure of the grid has been abstracted through an interface that merely provides 3 methods
expected of any grid-based implementation
*/
package fov

import (
	"math"
)

type GridMap interface {
	InBounds(x, y int) bool
	IsOpaque(x, y int) bool
}

type point struct {
	x, y int
}

type gridSet map[point]struct{}

// visible set of cells any time it is called
type View struct {
	Visible gridSet
}

func New() *View {
	return &View{}
}

func (v *View) Compute(grid GridMap, px, py, radius int) {
	v.Visible = make(map[point]struct{})
	v.Visible[point{px, py}] = struct{}{}
	for i := 1; i <= 8; i++ {
		v.fov(grid, px, py, 1, 0, 1, i, radius)
	}
}

func (v *View) fov(grid GridMap, px, py, dist int, lowSlope, highSlope float64, oct, rad int) {
	if dist > rad {
		return
	}

	low := math.Floor(lowSlope*float64(dist) + 0.5)
	high := math.Floor(highSlope*float64(dist) + 0.5)

	inGap := false

	for height := low; height <= high; height++ {
		mapx, mapy := distHeightXY(px, py, dist, int(height), oct)
		if grid.InBounds(mapx, mapy) && distanceTo(px, py, mapx, mapy) < rad {
			v.Visible[point{mapx, mapy}] = struct{}{}
		}

		if grid.InBounds(mapx, mapy) && grid.IsOpaque(mapx, mapy) {
			if inGap {
				v.fov(grid, px, py, dist+1, lowSlope, (height-0.5)/float64(dist), oct, rad)
			}
			lowSlope = (height + 0.5) / float64(dist)
			inGap = false
		} else {
			inGap = true
			if height == high {
				v.fov(grid, px, py, dist+1, lowSlope, highSlope, oct, rad)
			}
		}
	}
}

func (v *View) IsVisible(x, y int) bool {
	if _, ok := v.Visible[point{x, y}]; ok {
		return true
	}
	return false
}

func distHeightXY(px, py, d, h, oct int) (int, int) {
	if oct&0x1 > 0 {
		d = -d
	}
	if oct&0x2 > 0 {
		h = -h
	}
	if oct&0x4 > 0 {
		return px + h, py + d
	}
	return px + d, py + h
}

func distanceTo(x1, y1, x2, y2 int) int {
	vx := math.Pow(float64(x1-x2), 2)
	vy := math.Pow(float64(y1-y2), 2)
	return int(math.Sqrt(vx + vy))
}