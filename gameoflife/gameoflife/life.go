package gameoflife

type Cell struct {
    X, Y int64
}

type Cells []Cell

func NewGame(configuration Cells) *Game {
    cells := make(map[Cell]bool)
    for _, c := range configuration {
        cells[c] = true
    }

    return &Game{
        cells: cells,
    }
}

type Game struct {
    cells map[Cell]bool
    genNum int64
}

func (g *Game) NextGeneration() (Cells, int64) {
    neighborsCount := g.countNeighbors()
    g.cells = g.changeGeneration(neighborsCount)

    var cells Cells
    for c, v := range g.cells {
        if v {
            cells = append(cells, c)
        }
    }

    g.genNum++

    return cells, g.genNum
}

func (g *Game) countNeighbors() map[Cell]int {
    neighborsCount := make(map[Cell]int)
    for c, v := range g.cells {
        if !v {
            continue
        }

        if _, ok := neighborsCount[c]; !ok {
            neighborsCount[c] = 0
        }

        g.forEachNeighbor(c, func(n Cell) {
            neighborsCount[n]++
        })
    }

    return neighborsCount
}

func (g *Game) forEachNeighbor(c Cell, f func(n Cell)) {
    for x := int64(-1); x < 2; x++ {
        for y := int64(-1); y < 2; y++ {
            if x == 0 && y == 0 {
                continue
            }

            f(Cell{c.X + x, c.Y + y})
        }
    }
}

func (g *Game) changeGeneration(neighborsCount map[Cell]int) map[Cell]bool {
    gen := make(map[Cell]bool)
    for c, v := range neighborsCount {
        if v == 2 && g.cells[c] {
            gen[c] = true
        } else if v == 3 {
            gen[c] = true
        } else if v < 2 || v > 3 {
            gen[c] = false
        }
    }

    return gen
}
