package gameoflife

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestSingleCellDeath(t *testing.T) {
    configuration := Cells{{0, 0}}
    game := NewGame(configuration)

    cells, _ := game.NextGeneration()

    require.Empty(t, cells)
}

func TestTwoAdjacentCellsDeath(t *testing.T) {
    configuration := Cells{{0, 0}, {0, 1}}
    game := NewGame(configuration)

    cells, _ := game.NextGeneration()

    require.Empty(t, cells)
}

func TestTwoRemoteCellsDeath(t *testing.T) {
    configuration := Cells{{0, 0}, {100, 100}}
    game := NewGame(configuration)

    cells, _ := game.NextGeneration()

    require.Empty(t, cells)
}

func TestSurviveOneGeneration(t *testing.T) {
    configuration := Cells{{0, 0}, {1, 1}, {2, 1}}
    game := NewGame(configuration)

    cells, genNum := game.NextGeneration()

    require.Equal(t, 2, len(cells))
    require.Contains(t, cells, Cell{1, 0}, Cell{1, 1})
    require.Equal(t, int64(1), genNum)

    cells, genNum = game.NextGeneration()

    require.Empty(t, cells)
    require.Equal(t, int64(2), genNum)
}

func TestCountNeighbors(t *testing.T) {
    configuration := Cells{{1, 0}, {1, 1}}
    game := NewGame(configuration)

    neighboursCount := game.countNeighbors()

    require.Equal(t, map[Cell]int{
        Cell{0, 0}:  2,
        Cell{2, 0}:  2,
        Cell{0, 1}:  2,
        Cell{2, 1}:  2,
        Cell{2, 0}:  2,
        Cell{1, 0}:  1,
        Cell{1, 1}:  1,
        Cell{0, 2}:  1,
        Cell{1, 2}:  1,
        Cell{2, 2}:  1,
        Cell{0, -1}: 1,
        Cell{1, -1}: 1,
        Cell{2, -1}: 1,
    }, neighboursCount)
}

func TestCompleteConfiguration(t *testing.T) {
    configuration := Cells{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
    game := NewGame(configuration)

    cells, _ := game.NextGeneration()

    require.Equal(t, 4, len(cells))
    require.Contains(t, cells, Cell{0, 0})
    require.Contains(t, cells, Cell{1, 0})
    require.Contains(t, cells, Cell{0, 1})
    require.Contains(t, cells, Cell{1, 1})
}
