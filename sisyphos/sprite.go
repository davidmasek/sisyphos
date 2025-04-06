package sisyphos

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite represents an image.
type Sprite struct {
	image      *ebiten.Image
	alphaImage *image.Alpha
	x          int
	y          int
	action     func()
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	return image.Point{x, y}.In(s.image.Bounds())
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) InAlpha(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Use alphaImage (*image.Alpha) instead of image (*ebiten.Image) here.
	// It is because (*ebiten.Image).At is very slow as this reads pixels from GPU,
	// and should be avoided whenever possible.
	if s.alphaImage == nil {
		s.alphaImage = image.NewAlpha(s.image.Bounds())
		cloneToAlpha(s.image, s.alphaImage)
	}
	return s.alphaImage.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

func (s *Sprite) JustPressed() {
	s.action()
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(s.image, op)
}

func (g *Game) spriteAt(x, y int) *Sprite {
	// As the sprites are ordered from back to front,
	// search the clicked/touched sprite in reverse order.
	for i := len(g.sprites) - 1; i >= 0; i-- {
		s := g.sprites[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *Game) moveSpriteToFront(sprite *Sprite) {
	index := -1
	for i, ss := range g.sprites {
		if ss == sprite {
			index = i
			break
		}
	}
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
	g.sprites = append(g.sprites, sprite)
}
