package sisyphos

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/*
var assetsFolder embed.FS

var (
	tileImage     = ebiten.NewImage(tileSize, tileSize)
	playerImage   = ebiten.NewImage(tileSize, tileSize)
	boulderImage  = ebiten.NewImage(tileSize, tileSize)
	mountainImage = ebiten.NewImage(tileSize, tileSize)
	targetImage   = ebiten.NewImage(tileSize, tileSize)

	restartImage = ebiten.NewImage(tileSize, tileSize)

	mplusFaceSource *text.GoTextFaceSource
)

func cloneToAlpha(src *ebiten.Image, target *image.Alpha) {
	// Clone an image but only with alpha values.
	// This is used to detect a user cursor touches the image.
	b := src.Bounds()
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			target.Set(i, j, src.At(i, j))
		}
	}
}

func init() {
	tileImage.Fill(color.White)

	loadImage("assets/stickman.png", playerImage)
	loadImage("assets/boulder.png", boulderImage)
	loadImage("assets/mountain.png", mountainImage)
	loadImage("assets/vase.png", targetImage)

	loadImage("assets/restart.png", restartImage)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func loadImage(path string, target *ebiten.Image) {
	imgByte, err := assetsFolder.ReadFile(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileSize)/float64(w), float64(tileSize)/float64(h))

	// Draw the original decoded image onto 'resized'
	src := ebiten.NewImageFromImage(img)
	target.DrawImage(src, op)
}
