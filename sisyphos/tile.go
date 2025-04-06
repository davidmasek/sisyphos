package sisyphos

import (
	"errors"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

// TileData represents a tile information like a value and a position.
type TileData struct {
	value SpriteType
	x     int
	y     int
}

// Tile represents a tile information including TileData and animation states.
type Tile struct {
	current TileData

	// next represents a next tile information after moving.
	// next is empty when the tile is not about to move.
	next TileData

	movingCount       int
	startPoppingCount int
	poppingCount      int
}

// Pos returns the tile's current position.
// Pos is used only at testing so far.
func (t *Tile) Pos() (int, int) {
	return t.current.x, t.current.y
}

// NextPos returns the tile's next position.
// NextPos is used only at testing so far.
func (t *Tile) NextPos() (int, int) {
	return t.next.x, t.next.y
}

// Value returns the tile's current value.
// Value is used only at testing so far.
func (t *Tile) Value() SpriteType {
	return t.current.value
}

// NextValue returns the tile's current value.
// NextValue is used only at testing so far.
func (t *Tile) NextValue() SpriteType {
	return t.next.value
}

// NewTile creates a new Tile object.
func NewTile(value SpriteType, x, y int) *Tile {
	return &Tile{
		current: TileData{
			value: value,
			x:     x,
			y:     y,
		},
		startPoppingCount: maxPoppingCount,
	}
}

// IsMoving returns a boolean value indicating if the tile is animating.
func (t *Tile) IsMoving() bool {
	return 0 < t.movingCount
}

func (t *Tile) stopAnimation() {
	if 0 < t.movingCount {
		t.current = t.next
		t.next = TileData{}
	}
	t.movingCount = 0
	t.startPoppingCount = 0
	t.poppingCount = 0
}

func tileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
	var result *Tile
	for t := range tiles {
		if t.current.x != x || t.current.y != y {
			continue
		}
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

const ()

// MoveTiles moves tiles in the given tiles map if possible.
// MoveTiles returns true if there are tiles that are to move, otherwise false.
//
// When MoveTiles is called, all tiles must not be about to move.
func MoveTiles(tiles map[*Tile]struct{}, size int, dir Dir) bool {
	for t := range tiles {
		if t.current.value == PlayerSprite {
			nx, ny := t.current.x, t.current.y
			nnx, nny := nx, ny
			switch dir {
			case DirUp:
				ny -= 1
				nny -= 2
			case DirDown:
				ny += 1
				nny += 2
			case DirLeft:
				nx -= 1
				nnx -= 2
			case DirRight:
				nx += 1
				nnx += 2
			}
			next := tileAt(tiles, nx, ny)
			if next == nil || next.Value() == TargetSprite {
				nextData := TileData{t.current.value, nx, ny}
				t.next = nextData
				t.movingCount = maxMovingCount
				return true
			}
			if next.Value() == BoulderSprite {
				nNext := tileAt(tiles, nnx, nny)
				if nNext == nil || nNext.Value() == TargetSprite {
					nextData := TileData{t.current.value, nx, ny}
					t.next = nextData
					t.movingCount = maxMovingCount

					nextNextData := TileData{next.Value(), nnx, nny}
					next.next = nextNextData
					next.movingCount = maxMovingCount
					return true
				}

			}
		}
	}
	return false
}

func addRandomTile(tiles map[*Tile]struct{}, size int, sprite SpriteType) error {
	cells := make([]bool, size*size)
	for t := range tiles {
		if t.IsMoving() {
			panic("not reach")
		}
		i := t.current.x + t.current.y*size
		cells[i] = true
	}
	availableCells := []int{}
	for i, b := range cells {
		if b {
			continue
		}
		availableCells = append(availableCells, i)
	}
	if len(availableCells) == 0 {
		return errors.New("sisyphos: there is no space to add a new tile")
	}
	c := availableCells[rand.IntN(len(availableCells))]
	x := c % size
	y := c / size
	t := NewTile(sprite, x, y)
	tiles[t] = struct{}{}
	return nil
}

// Update updates the tile's animation states.
func (t *Tile) Update() error {
	switch {
	case 0 < t.movingCount:
		t.movingCount--
		if t.movingCount == 0 {
			if t.current.value != t.next.value && 0 < t.next.value {
				t.poppingCount = maxPoppingCount
			}
			t.current = t.next
			t.next = TileData{}
		}
	case 0 < t.startPoppingCount:
		t.startPoppingCount--
	case 0 < t.poppingCount:
		t.poppingCount--
	}
	return nil
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}

// Draw draws the current tile to the given boardImage.
func (t *Tile) Draw(boardImage *ebiten.Image) {
	i, j := t.current.x, t.current.y
	ni, nj := t.next.x, t.next.y
	v := t.current.value
	if v == EmptySprite {
		return
	}
	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	nx := ni*tileSize + (ni+1)*tileMargin
	ny := nj*tileSize + (nj+1)*tileMargin
	switch {
	case 0 < t.movingCount:
		rate := 1 - float64(t.movingCount)/maxMovingCount
		x = mean(x, nx, rate)
		y = mean(y, ny, rate)
	case 0 < t.startPoppingCount:
		rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
		scale := meanF(0.0, 1.0, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	case 0 < t.poppingCount:
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			// 0 to 1
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			// 1 to 0
			rate = float64(t.poppingCount) / float64(maxPoppingCount*2/3)
		}
		scale := meanF(1.0, maxScale, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	}
	op.GeoM.Translate(float64(x), float64(y))
	boardImage.DrawImage(tileSprite(v), op)
}

func tileSprite(value SpriteType) *ebiten.Image {
	switch value {
	case PlayerSprite:
		return playerImage
	case BoulderSprite:
		return boulderImage
	case MountainSprite:
		return mountainImage
	case TargetSprite:
		return targetImage
	}
	log.Println(value)
	panic("not reach")
}
