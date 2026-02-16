# Ebiten Project

A small turn-based dungeon crawler built with Ebiten. It generates a dungeon, places the player and monsters, and runs a simple combat loop with a HUD and message log.

## Features

- Procedural dungeon rooms and corridors
- Turn-based movement and combat
- Field of view (FOV) and fog-of-war exploration
- Simple AI pathfinding (A*) for monsters
- HUD with player stats and a scrolling combat log

## Controls

- Move: Arrow keys or WASD
- Skip turn: Q

## Download (binaries)

- Latest release: https://github.com/Sanjar0126/ebiten_project/releases/tag/v0.0.1
- All versions: https://github.com/Sanjar0126/ebiten_project/tags

Releases include binaries for Linux and Windows (.exe).
The assets folder is required to run the binaries.

## Run from source

### Requirements

- Go 1.22.1 or newer

### Steps

```bash
go mod download
go run .
```

### Build

```bash
go build -o ebiten_project
```

On Windows:

```bash
go build -o ebiten_project.exe
```

## Project layout

- Core game loop: main.go
- World setup and ECS entities: world.go, components.go
- Dungeon and level generation: level.go, map.go, dungeon.go
- Combat and turn state: combat.go, turnstate.go
- Rendering and UI: render.go, hud.go, userlog.go
- Pathfinding: path_finding.go
- Assets: assets/

## Notes

This project uses Ebiten for rendering and input, and Arche ECS for entity management.
