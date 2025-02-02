package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
)

type Game struct {
	Map         GameMap
	World       *ecs.World
	WorldTags   map[string]ecs.Entity
	Turn        TurnState
	TurnCounter int
}

func (g *Game) Update() error {
	g.TurnCounter++
	if g.Turn == PlayerTurn && g.TurnCounter > 20 {
		TakePlayerAction(g)
	}

	if g.Turn == MonsterTurn {
		UpdateMonster(g)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	level := g.Map.CurrentLevel
	level.DrawLevel(screen)
	ProcessRenderables(g, level, screen)
	ProcessUserLog(g, screen)
	ProcessHUD(g, screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %v | FPS: %v\n", ebiten.ActualTPS(), ebiten.ActualFPS()))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Turn: %v | Moves: %v\n", g.Turn, g.TurnCounter))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	gd := NewGameData()
	return gd.TileWidth * gd.ScreenWidth, gd.TileHeight * gd.ScreenHeight
}

func NewGame() *Game {
	g := &Game{}
	g.Map = NewGameMap()
	world := InitializeWorld(g.Map.CurrentLevel)
	g.World = world
	g.Turn = PlayerTurn
	g.TurnCounter = 0
	return g
}

func main() {
	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
