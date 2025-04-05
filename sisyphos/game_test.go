package sisyphos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
	const size = 4
	game, err := NewGame()
	require.NoError(t, err)
	for tile := range game.board.tiles {
		t.Logf("%#v\n", tile)
	}
}
