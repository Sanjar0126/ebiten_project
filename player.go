package main

import "github.com/hajimehoshi/ebiten/v2"

func TakePlayerAction(g *Game) {
	players := playerFilter.Query(g.World)
	turnTaken := false

	x := 0
	y := 0
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		x = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		turnTaken = true
	}

	level := g.Map.CurrentLevel

	for players.Next() {
		_, pos, _, _, _, _, _, _, _, _ := players.Get()
		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)

		tile := level.Tiles[index]
		if !tile.Blocked {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false

			pos.X += x
			pos.Y += y
			level.Tiles[index].Blocked = true
			level.PlayerFov.Compute(level, pos.X, pos.Y, 8)
		} else if x != 0 || y != 0 {
			if level.Tiles[index].TileType != WALL {
				monsterPosition := Position{X: pos.X + x, Y: pos.Y + y}
				AttackSystem(g, pos, &monsterPosition)
			}
		}
	}

	if x != 0 || y != 0 || turnTaken {
		g.Turn = GetNextState(g.Turn)
		g.TurnCounter = 0
	}
}
