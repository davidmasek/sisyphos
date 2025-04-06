package sisyphos

import (
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth    = 600
	ScreenHeight   = 600
	StartBoardSize = 3
	startBlocks    = 0
	StartX         = 1
	StartY         = 1

	tileSize   = 90
	tileMargin = 4

	// controls movement speed
	maxMovingCount  = 5
	maxPoppingCount = 6

	MinDragDistance = 8
)

type SpriteType int

const (
	EmptySprite SpriteType = iota
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
	boardSize  int

	sprites []*Sprite
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		input:     NewInput(),
		level:     0,
		boardSize: StartBoardSize,
	}
	g.restart()

	// Initialize the sprites.
	sprites := []*Sprite{}
	// w, h := restartImage.Bounds().Dx(), restartImage.Bounds().Dy()
	restart := &Sprite{
		image: restartImage,
		x:     0,
		y:     0,
		action: func() {
			log.Println("restart button pressed")
			g.restart()
		},
	}
	sprites = append(sprites, restart)

	g.sprites = sprites

	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) restart() {
	var err error
	retries := 0
	g.board, err = NewBoard(g.boardSize, startBlocks+g.level)
	for err != nil {
		g.boardSize += 1
		g.boardImage = nil
		g.board, err = NewBoard(g.boardSize, startBlocks+g.level)
		// safeguard in case we can never generate the game
		if retries > 100 {
			panic("cannot restart game")
		}
		retries += 1
	}
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	for _, pos := range g.input.Clicks {
		startSprite := g.spriteAt(pos.StartX, pos.StartY)
		endSprite := g.spriteAt(pos.EndX, pos.EndY)
		if startSprite != nil && startSprite == endSprite {
			g.moveSpriteToFront(endSprite)
			endSprite.JustPressed()
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.restart()
	}
	if gameOver(g.board) || inpututil.IsKeyJustReleased(ebiten.KeyU) {
		g.level += 1
		g.restart()
	}
	if runtime.GOOS != "js" && inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		return ebiten.Termination
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

	for _, s := range g.sprites {
		s.Draw(screen, 1)
	}
}
