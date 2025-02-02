package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(g *Game, level Level, screen *ebiten.Image) {
	query := renderableFilter.Query(g.World)
	for query.Next() {
		pos, ren := query.Get()
		img := ren.Image
		isVisible := level.PlayerFov.IsVisible(pos.X, pos.Y)

		if isVisible {
			index := level.GetIndexFromXY(pos.X, pos.Y)
			tile := level.Tiles[index]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(img, op)
		}
	}
}
