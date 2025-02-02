package main

import (
	"fmt"
	"log"

	"github.com/Sanjar0126/ebiten_project/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var hudImg *ebiten.Image = nil
var hudErr error = nil
var hudFont font.Face = nil

func ProcessHUD(g *Game, screen *ebiten.Image) {
	if hudImg == nil {
		hudImg, _, hudErr = ebitenutil.NewImageFromFile("assets/UIPanel.png")
		if hudErr != nil {
			log.Fatal(hudErr)
		}
	}
	if hudFont == nil {
		tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}

		const dpi = 72
		hudFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    16,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	gd := NewGameData()

	uiY := (gd.ScreenHeight - gd.UIHeight) * gd.TileHeight
	uiX := (gd.ScreenWidth * gd.TileWidth) / 2
	var fontX = uiX + 16
	var fontY = uiY + 24
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(uiX), float64(uiY))
	screen.DrawImage(userLogImg, op)

	playerQuery := playerFilter.Query(g.World)
	for playerQuery.Next() {
		_, _, _, _, health, _, weapon, armor, _, _ := playerQuery.Get()

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(fontX), float64(fontY))
		op.LineSpacing = mplusNormalFont.Size * 1.5

		healthText := fmt.Sprintf("Health: %d / %d", health.CurrentHealth, health.MaxHealth)
		text.Draw(screen, healthText, mplusNormalFont, op)

		op = &text.DrawOptions{}
		fontY += 16
		op.GeoM.Translate(float64(fontX), float64(fontY))

		acText := fmt.Sprintf("Armor Class: %s", ArmorClassNameMap[armor.ArmorClass])
		text.Draw(screen, acText, mplusNormalFont, op)

		op = &text.DrawOptions{}
		fontY += 16
		op.GeoM.Translate(float64(fontX), float64(fontY))
		
		defText := fmt.Sprintf("Defense: %d", armor.Defense)
		text.Draw(screen, defText, mplusNormalFont, op)

		op = &text.DrawOptions{}
		fontY += 16
		op.GeoM.Translate(float64(fontX), float64(fontY))

		dmg := fmt.Sprintf("Damage: %d - %d", weapon.MinimumDamage, weapon.MaximumDamage)
		text.Draw(screen, dmg, mplusNormalFont, op)
	}
}
