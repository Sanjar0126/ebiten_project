package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type UserMessage struct {
	AttackMessage    string
	DeadMessage      string
	GameStateMessage string
}

type WeaponType int

const (
	Melee = iota
	Ranged
	Magic
)

type ArmorClass int

const (
	Light = iota
	Medium
	Heavy
)

var (
	ArmorClassNameMap = map[int]string{
		Light:  "Light",
		Medium: "Medium",
		Heavy:  "Heavy",
	}
)

const (
	SkeletonName = "skeleton"
	ZombieName   = "zombie"
)

var (
	MonsterNameMap = map[string]string{
		SkeletonName: "Skeleton",
		ZombieName:   "Zombie",
	}

	MonsterFovRangeMap = map[string]int{
		SkeletonName: 8,
		ZombieName:   8,
	}
)

type Player struct{}

type Health struct {
	MaxHealth     int
	CurrentHealth int
}

type Mana struct {
	MaxMana     int
	CurrentMana int
}

type Stats struct {
}

type Weapon struct {
	Type          WeaponType
	Name          string
	MinimumDamage int
	MaximumDamage int
	ToHitBonus    int
	CritChance    int
	CritMultip    int
	ArmorPierce   int
}

type Armor struct {
	Name        string
	Defense     int
	EvadeChance int
	ArmorClass  int
}

type Position struct {
	X int
	Y int
}

type Renderable struct {
	Image *ebiten.Image
}

type Movable struct{}

type Monster struct{}

type Name struct {
	Label string
}

func (p *Position) GetManhattanDistance(other *Position) int {
	xDist := math.Abs(float64(p.X - other.X))
	yDist := math.Abs(float64(p.Y - other.Y))
	return int(xDist) + int(yDist)
}

func (p *Position) IsEqual(other *Position) bool {
	return (p.X == other.X && p.Y == other.Y)
}
