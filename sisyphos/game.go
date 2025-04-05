package sisyphos

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
	boardSize    = 5
	startBlocks  = 2
	StartX       = 2
	StartY       = 2
)

type Sprite int

const (
	EmptySprite Sprite = iota
	PlayerSprite
	BoulderSprite
	MountainSprite
	TargetSprite
)

// Game represents a game state.
type Game struct {
	input      *Input
	board      *Board
	boardImage *ebiten.Image
	level      int
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
		level: 0,
	}
	var err error
	g.board, err = NewBoard(boardSize, startBlocks)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	if gameOver(g.board) {
		g.level += 1
		var err error
		g.board, err = NewBoard(boardSize, startBlocks+g.level)
		if err != nil {
			return err
		}
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(g.board.Size())
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
}
