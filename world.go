package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

var (
	playerBuilder  generic.Map10[Player, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]
	monsterBuilder generic.Map10[Monster, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]

	renderableFilter *generic.Filter2[Position, Renderable]
	playerFilter     *generic.Filter10[Player, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]
	monsterFilter    *generic.Filter10[Monster, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]
	messangerFilter  *generic.Filter2[UserMessage, Name]
)

func InitializeWorld(startingLevel Level) *ecs.World {
	world := ecs.NewWorld()

	startingRoom := startingLevel.Rooms[0]
	x, y := startingRoom.Center()

	initPlayer(&world, x, y, getStats(statsArgs{
		maxHealth:           30,
		currentHealth:       30,
		weaponType:          Melee,
		weaponName:          "Long Sword",
		weaponMinimumDamage: 10,
		weaponMaximumDamage: 20,
		weaponToHitBonus:    0,
		weaponCritChance:    0,
		weaponCritMultip:    0,
		weaponArmorPierce:   0,
		armorName:           "Knight armor",
		armorDefense:        15,
		armorArmorClass:     Heavy,
		label:               "Player",
	}))

	skeletonStats := getStats(statsArgs{
		maxHealth:           10,
		currentHealth:       10,
		weaponType:          Melee,
		weaponName:          "Rusted Short Sword",
		weaponMinimumDamage: 2,
		weaponMaximumDamage: 6,
		weaponToHitBonus:    0,
		weaponCritChance:    0,
		weaponCritMultip:    0,
		weaponArmorPierce:   0,
		armorName:           "Bone",
		armorDefense:        2,
		armorArmorClass:     Light,
		label:               SkeletonName,
	})
	zombieStats := getStats(statsArgs{
		maxHealth:           20,
		currentHealth:       20,
		weaponType:          Melee,
		weaponName:          "Rusted axe",
		weaponMinimumDamage: 5,
		weaponMaximumDamage: 10,
		weaponToHitBonus:    0,
		weaponCritChance:    0,
		weaponCritMultip:    0,
		weaponArmorPierce:   0,
		armorName:           "Rusted chainmail",
		armorDefense:        8,
		armorArmorClass:     Medium,
		label:               ZombieName,
	})

	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()
			mobSpawn := GetDiceRoll(2)
			var monsterType string
			var monsterStats stats
			if mobSpawn == 1 {
				monsterType = SkeletonName
				monsterStats = skeletonStats
			} else {
				monsterType = ZombieName
				monsterStats = zombieStats
			}

			initMonster(&world, mX, mY, monsterType, monsterStats)
		}

	}

	renderableFilter = generic.NewFilter2[Position, Renderable]()
	playerFilter = generic.NewFilter10[Player, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]()
	monsterFilter = generic.NewFilter10[Monster, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage]()
	messangerFilter = generic.NewFilter2[UserMessage, Name]()

	return &world
}

func initPlayer(world *ecs.World, x, y int, stats stats) {
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	playerBuilder = generic.NewMap10[Player, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage](world)

	_ = playerBuilder.NewWith(&Player{}, &Position{X: x, Y: y}, &Renderable{Image: playerImg}, &Movable{},
		&stats.health, &stats.mana, &stats.weapon, &stats.armor, &stats.name, &UserMessage{})

}

func initMonster(world *ecs.World, x, y int, monsterName string, stats stats) {
	monsterImg, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("assets/%s.png", monsterName))
	if err != nil {
		log.Fatal(err)
	}

	monsterBuilder = generic.NewMap10[Monster, Position, Renderable, Movable, Health, Mana, Weapon, Armor, Name, UserMessage](world)

	_ = monsterBuilder.NewWith(&Monster{}, &Position{X: x, Y: y}, &Renderable{Image: monsterImg}, &Movable{},
		&stats.health, &stats.mana, &stats.weapon, &stats.armor, &stats.name, &UserMessage{})
}

type stats struct {
	health Health
	mana   Mana
	weapon Weapon
	armor  Armor
	name   Name
}

type statsArgs struct {
	maxHealth     int
	currentHealth int

	maxMana     int
	currentMana int

	weaponType          WeaponType
	weaponName          string
	weaponMinimumDamage int
	weaponMaximumDamage int
	weaponToHitBonus    int
	weaponCritChance    int
	weaponCritMultip    int
	weaponArmorPierce   int

	armorName       string
	armorDefense    int
	armorArmorClass int

	label string
}

func getStats(args statsArgs) stats {
	return stats{
		health: Health{
			MaxHealth:     args.maxHealth,
			CurrentHealth: args.currentHealth,
		},
		mana: Mana{
			MaxMana:     args.maxMana,
			CurrentMana: args.currentMana,
		},
		weapon: Weapon{
			Type:          args.weaponType,
			Name:          args.weaponName,
			MinimumDamage: args.weaponMinimumDamage,
			MaximumDamage: args.weaponMaximumDamage,
			ToHitBonus:    args.weaponToHitBonus,
			CritChance:    args.weaponCritChance,
			CritMultip:    args.weaponCritMultip,
			ArmorPierce:   args.weaponArmorPierce,
		},
		armor: Armor{
			Name:       args.armorName,
			Defense:    args.armorDefense,
			ArmorClass: args.armorArmorClass,
		},
		name: Name{
			Label: args.label,
		},
	}
}
