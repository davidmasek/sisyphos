package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"sisyphos.optimisticotter.me/sisyphos"
)

func main() {
	game, err := sisyphos.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(sisyphos.ScreenWidth, sisyphos.ScreenHeight)
	ebiten.SetWindowTitle("sisyphos.optimisticotter.me")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
