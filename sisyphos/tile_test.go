package sisyphos_test

import (
	"fmt"
	"testing"

	"sisyphos.optimisticotter.me/sisyphos"
)

func cellsToTiles(cells []sisyphos.SpriteType, size int) map[*sisyphos.Tile]struct{} {
	tiles := map[*sisyphos.Tile]struct{}{}
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			c := cells[i+j*size]
			if c == 0 {
				continue
			}
			t := sisyphos.NewTile(c, i, j)
			tiles[t] = struct{}{}
		}
	}
	return tiles
}

func tilesToCells(tiles map[*sisyphos.Tile]struct{}, size int) ([]sisyphos.SpriteType, []sisyphos.SpriteType) {
	cells := make([]sisyphos.SpriteType, size*size)
	nextCells := make([]sisyphos.SpriteType, size*size)
	for t := range tiles {
		x, y := t.Pos()
		cells[x+y*size] = t.Value()
		if t.IsMoving() {
			if t.NextValue() == 0 {
				continue
			}
			nx, ny := t.NextPos()
			nextCells[nx+ny*size] = t.NextValue()
		} else {
			nextCells[x+y*size] = t.Value()
		}
	}
	return cells, nextCells
}

func TestMoveTiles(t *testing.T) {
	const size = 4
	testCases := []struct {
		Dir   sisyphos.Dir
		Input []sisyphos.SpriteType
		Want  []sisyphos.SpriteType
	}{
		{
			Dir: sisyphos.DirUp,
			Input: []sisyphos.SpriteType{
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
			},
			Want: []sisyphos.SpriteType{
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
				sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite, sisyphos.EmptySprite,
			},
		},
	}
	for _, test := range testCases {
		want, _ := tilesToCells(cellsToTiles(test.Want, size), size)
		tiles := cellsToTiles(test.Input, size)
		moved := sisyphos.MoveTiles(tiles, size, test.Dir)
		input, got := tilesToCells(tiles, size)
		if !moved {
			got = input
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("dir: %s, input: %v, got %v; want %v", test.Dir.String(), test.Input, got, want)
		}
	}
}
