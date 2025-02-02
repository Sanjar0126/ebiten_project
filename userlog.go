package main

import (
	"bytes"
	"log"

	"github.com/Sanjar0126/ebiten_project/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mlange-42/arche/ecs"
)

var userLogImg *ebiten.Image = nil
var err error = nil
var mplusNormalFont *text.GoTextFace = nil
var mplusFaceSource *text.GoTextFaceSource
var lastText []string = make([]string, 0, 5)
var toRemove []ecs.Entity

func ProcessUserLog(g *Game, screen *ebiten.Image) {
	if userLogImg == nil {
		userLogImg, _, err = ebitenutil.NewImageFromFile("assets/UIPanel.png")
		if err != nil {
			log.Fatal(err)
		}
	}
	if mplusNormalFont == nil {
		s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
		if err != nil {
			log.Fatal(err)
		}
		mplusFaceSource = s

		mplusNormalFont = &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   12,
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	gd := NewGameData()

	uiLocation := (gd.ScreenHeight - gd.UIHeight) * gd.TileHeight
	var fontX = 16
	var fontY = uiLocation + 24
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0.), float64(uiLocation))
	screen.DrawImage(userLogImg, op)
	tmpMessages := make([]string, 0, 5)
	anyMessages := false

	messangerQuery := messangerFilter.Query(g.World)
	for messangerQuery.Next() {
		messages, name := messangerQuery.Get()
		if messages.AttackMessage != "" {
			tmpMessages = append(tmpMessages, messages.AttackMessage)
			anyMessages = true
			messages.AttackMessage = ""
		} else if messages.DeadMessage != "" {
			if messages.DeadMessage != "" {
				tmpMessages = append(tmpMessages, messages.DeadMessage)
				anyMessages = true
				messages.DeadMessage = ""
				if name.Label != "Player" {
					toRemove = append(toRemove, messangerQuery.Entity())
				}
			}
			if messages.GameStateMessage != "" {
				tmpMessages = append(tmpMessages, messages.GameStateMessage)
				anyMessages = true
			}
		}
	}

	messangerQuery = messangerFilter.Query(g.World)
	for messangerQuery.Next() {
		messages, name := messangerQuery.Get()
		if messages.DeadMessage != "" {
			tmpMessages = append(tmpMessages, messages.DeadMessage)
			anyMessages = true
			messages.DeadMessage = ""
			if name.Label != "Player" {
				toRemove = append(toRemove, messangerQuery.Entity())
			}
		}

		if messages.GameStateMessage != "" {
			tmpMessages = append(tmpMessages, messages.GameStateMessage)
			anyMessages = true
		}
	}

	for _, e := range toRemove {
		g.World.RemoveEntity(e)
	}

	if anyMessages {
		lastText = append(tmpMessages, lastText...)
	}
	if len(lastText) > 8 {
		lastText = lastText[:8]
	}
	for _, msg := range lastText {
		if msg != "" {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(fontX), float64(fontY))
			op.LineSpacing = mplusNormalFont.Size * 1.5
			text.Draw(screen, msg, mplusNormalFont, op)
			fontY += 16
		}
	}

}
