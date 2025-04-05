package sisyphos

import (
	"image/color"
	"log"
)

var (
	backgroundColor = color.RGBA{75, 75, 75, 0xff}
	frameColor      = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
)

func tileBackgroundColor(value Sprite) color.Color {
	switch value {
	case EmptySprite:
		return color.NRGBA{0xee, 0xe4, 0xda, 0x59}
	case PlayerSprite:
		return color.RGBA{0xee, 0xe4, 0xda, 0xff}
	case BoulderSprite:
		return color.RGBA{0xed, 0xe0, 0xc8, 0xff}
	case MountainSprite:
		return color.RGBA{0xf2, 0xb1, 0x79, 0xff}
	}
	log.Println(value)
	panic("not reach")
}
