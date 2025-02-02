package main

import (
	"log"

	"github.com/Sanjar0126/ebiten_project/fov"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	floorImage *ebiten.Image
	wallImage  *ebiten.Image

	levelHeight int = 0
)

type TileType int

const (
	WALL TileType = iota
	FLOOR
)

type Level struct {
	Tiles     []*MapTile
	Rooms     []Rect
	PlayerFov *fov.View
}

func NewLevel() Level {
	loadTileImages()
	l := Level{}
	rooms := make([]Rect, 0)
	l.Rooms = rooms
	l.GenerateLevelTiles()
	l.PlayerFov = fov.New()
	return l
}

type MapTile struct {
	PixelX   int
	PixelY   int
	Blocked  bool
	Image    *ebiten.Image
	Explored bool
	TileType TileType
}

func (level *Level) GetIndexFromXY(x, y int) int {
	gd := NewGameData()
	levelHeight = gd.ScreenHeight - gd.UIHeight
	return (y * gd.ScreenWidth) + x
}

func (level *Level) GenerateLevelTiles() {
	MIN_SIZE := 6
	MAX_SIZE := 10
	MAX_ROOMS := 30

	gd := NewGameData()
	levelHeight = gd.ScreenHeight - gd.UIHeight
	tiles := level.createTiles()
	level.Tiles = tiles
	contains_rooms := false

	for idx := 0; idx < MAX_ROOMS; idx++ {
		w := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		h := GetRandomBetween(MIN_SIZE, MAX_SIZE)
		x := GetDiceRoll(gd.ScreenWidth - w - 1)
		y := GetDiceRoll(levelHeight - h - 1)
		new_room := NewRect(x, y, w, h)

		okToAdd := true

		for _, otherRoom := range level.Rooms {
			if new_room.Intersect(otherRoom) {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			level.createRoom(new_room)
			if contains_rooms {
				newX, newY := new_room.Center()
				prevX, prevY := level.Rooms[len(level.Rooms)-1].Center()
				coinflip := GetDiceRoll(2)
				if coinflip == 2 {
					level.createHorizontalTunnel(prevX, newX, prevY)
					level.createVerticalTunnel(prevY, newY, newX)
				} else {
					level.createHorizontalTunnel(prevX, newX, newY)
					level.createVerticalTunnel(prevY, newY, prevX)
				}
			}
			level.Rooms = append(level.Rooms, new_room)
			contains_rooms = true
		}
	}

}
func (level *Level) createTiles() []*MapTile {
	gd := NewGameData()
	tiles := make([]*MapTile, levelHeight*gd.ScreenWidth)
	index := 0

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
			index = level.GetIndexFromXY(x, y)

			tile := MapTile{
				PixelX:   x * gd.TileWidth,
				PixelY:   y * gd.TileHeight,
				Blocked:  true,
				Image:    wallImage,
				Explored: false,
			}
			tiles[index] = &tile
		}
	}

	return tiles
}

func (level *Level) createRoom(room Rect) {
	for y := room.Y1 + 1; y < room.Y2; y++ {
		for x := room.X1 + 1; x < room.X2; x++ {
			index := level.GetIndexFromXY(x, y)
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floorImage
		}
	}
}

func (level *Level) createHorizontalTunnel(x1 int, x2 int, y int) {
	gd := NewGameData()
	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := level.GetIndexFromXY(x, y)
		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floorImage
		}
	}
}

func (level *Level) createVerticalTunnel(y1 int, y2 int, x int) {
	gd := NewGameData()
	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := level.GetIndexFromXY(x, y)

		if index > 0 && index < gd.ScreenWidth*levelHeight {
			level.Tiles[index].Blocked = false
			level.Tiles[index].TileType = FLOOR
			level.Tiles[index].Image = floorImage
		}
	}
}

func (level *Level) DrawLevel(screen *ebiten.Image) {
	gd := NewGameData()
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < levelHeight; y++ {
			tileID := level.GetIndexFromXY(x, y)
			tile := level.Tiles[tileID]
			isVisible := level.PlayerFov.IsVisible(x, y)
			if isVisible {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				screen.DrawImage(tile.Image, op)
				level.Tiles[tileID].Explored = true
			}
			if tile.Explored && !isVisible {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				op.ColorScale.Scale(0.5, 0.5, 0.5, 1.0)
				screen.DrawImage(tile.Image, op)
			}
		}
	}
}

func (level Level) InBounds(x, y int) bool {
	gd := NewGameData()
	return x >= 0 && x <= gd.ScreenWidth && y >= 0 && y <= levelHeight
}

func (level Level) IsOpaque(x, y int) bool {
	index := level.GetIndexFromXY(x, y)
	return level.Tiles[index].TileType == WALL
}

func loadTileImages() {
	if floorImage != nil && wallImage != nil {
		return
	}
	var err error

	floorImage, _, err = ebitenutil.NewImageFromFile("assets/dirt.png")
	if err != nil {
		log.Fatal(err)
	}

	wallImage, _, err = ebitenutil.NewImageFromFile("assets/stone.png")
	if err != nil {
		log.Fatal(err)
	}
}
