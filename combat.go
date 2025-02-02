package main

import (
	"fmt"

	"github.com/mlange-42/arche/ecs"
)

var ToRemove []ecs.Entity

type attackerStats struct {
	position *Position
	health   *Health
	mana     *Mana
	weapon   *Weapon
	armor    *Armor
	name     *Name
	message  *UserMessage
}

func AttackSystem(g *Game, attackerPosition *Position, defenderPosition *Position) {
	var attacker *attackerStats
	var defender *attackerStats

	playerQuery := playerFilter.Query(g.World)
	for playerQuery.Next() {
		_, pos, _, _, health, mana, weapon, armor, name, message := playerQuery.Get()
		playerStats := attackerStats{
			position: pos,
			health:   health,
			mana:     mana,
			weapon:   weapon,
			armor:    armor,
			name:     name,
			message:  message,
		}

		if pos.IsEqual(attackerPosition) {
			attacker = &playerStats
		} else if pos.IsEqual(defenderPosition) {
			defender = &playerStats
		}
	}

	monsterQuery := monsterFilter.Query(g.World)
	for monsterQuery.Next() {
		_, pos, _, _, health, mana, weapon, armor, name, message := monsterQuery.Get()
		enemyStats := attackerStats{
			position: pos,
			health:   health,
			mana:     mana,
			weapon:   weapon,
			armor:    armor,
			name:     name,
			message:  message,
		}

		if pos.IsEqual(attackerPosition) {
			attacker = &enemyStats
			monsterQuery.Close()
			break
		} else if pos.IsEqual(defenderPosition) {
			defender = &enemyStats
			monsterQuery.Close()
			break
		}
	}

	if attacker == nil || defender == nil {
		return
	}

	if attacker.health.CurrentHealth <= 0 || defender.health.CurrentHealth <= 0 {
		return
	}

	damageRoll := GetRandomBetween(attacker.weapon.MinimumDamage, attacker.weapon.MaximumDamage)
	damageDone := damageRoll - defender.armor.Defense
	if damageDone < 0 {
		damageDone = 1
	}

	evadeRoll := GetRandomBetween(1, 100)
	if evadeRoll <= defender.armor.EvadeChance {
		damageDone = 0
		attacker.message.AttackMessage = "Attack miss!\n"
	}

	defender.health.CurrentHealth -= damageDone
	healthMsg := defender.health.CurrentHealth
	if healthMsg < 0 {
		healthMsg = 0
	}
	attacker.message.AttackMessage = fmt.Sprintf("%s swings %s at %s and hits for %d health, %s health: %d/%d\n", attacker.name.Label,
		attacker.weapon.Name, defender.name.Label, damageDone, defender.name.Label, healthMsg, defender.health.MaxHealth)

	if defender.health.CurrentHealth <= 0 {
		defender.message.DeadMessage = fmt.Sprintf("%s has died!\n", defender.name.Label)
		if defender.name.Label == "Player" {
			defender.message.GameStateMessage = "Game Over!\n"
			g.Turn = GameOver
		}
	}
}
