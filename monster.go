package main

import (
	"github.com/Sanjar0126/ebiten_project/fov"
	"github.com/mlange-42/arche/ecs"
)

func UpdateMonster(game *Game) {
	level := game.Map.CurrentLevel
	playerPos := Position{}

	playerQuery := playerFilter.Query(game.World)
	monsterQuery := monsterFilter.Query(game.World)

	for playerQuery.Next() {
		_, pos, _, _, _, _, _, _, _, _ := playerQuery.Get()
		playerPos.X = pos.X
		playerPos.Y = pos.Y
	}

	for monsterQuery.Next() {
		_, pos, _, _, health, _, _, _, name, _ := monsterQuery.Get()
		monsterFov := fov.New()
		monsterFov.Compute(level, pos.X, pos.Y, MonsterFovRangeMap[name.Label])
		if monsterFov.IsVisible(playerPos.X, playerPos.Y) {
			if pos.GetManhattanDistance(&playerPos) == 1 {
				AttackSystem(game, pos, &playerPos)
				if health.CurrentHealth <= 0 {
					t := level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)]
					t.Blocked = false
				}
			} else {
				aStar := AStar{}
				path := aStar.GetPath(level, pos, &playerPos)
				if len(path) > 0 {
					nextTile := level.Tiles[level.GetIndexFromXY(path[1].X, path[1].Y)]
					if !nextTile.Blocked {
						level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
						pos.X = path[1].X
						pos.Y = path[1].Y
						nextTile.Blocked = true
					}
				}
			}
		}
	}

	removeDeadEntities(game)

	if game.Turn != GameOver {
		game.Turn = PlayerTurn
	}
}

func removeDeadEntities(game *Game) {
	var toRemove []ecs.Entity
	monsterQuery := monsterFilter.Query(game.World)
	for monsterQuery.Next() {
		_, _, _, _, health, _, _, _, _, _ := monsterQuery.Get()
		if health.CurrentHealth <= 0 {
			toRemove = append(toRemove, monsterQuery.Entity())
		}
	}

	for _, entity := range toRemove {
		game.World.RemoveEntity(entity)
	}
}
