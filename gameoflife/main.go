package main

import (
    "fmt"
    "gameoflife/gameoflife"
    "os"
    "os/exec"
    "runtime"
    "time"
)

func main() {
    configuration := makeConfiguration()
    game := gameoflife.NewGame(configuration)

    m := makeMap(configuration)
    renderMap(m, 0)
    fmt.Println("Press Enter to start")
    _, _ = fmt.Scanln()

    cells := configuration
    var genNum int64
    for {
       m := makeMap(cells)
       renderMap(m, genNum)
       cells, genNum = game.NextGeneration()
       time.Sleep(delayBetweenGenerations)
    }
}

func renderMap(m [][]bool, genNum int64) {
    clear()
    fmt.Println("Generation: ", genNum)
    for _, row := range m {
        for _, c := range row {
            if c {
                fmt.Print("0")
            } else {
                fmt.Print(" ")
            }
        }
        fmt.Println()
    }
}

func makeMap(cells gameoflife.Cells) [][]bool {
    m := make([][]bool, size)
    for i := 0; i < size; i++ {
        m[i] = make([]bool, size)
    }

    for _, c := range cells {
        if c.Y >= 0 && c.Y < size && c.X >= 0 && c.X < size {
            m[c.Y][c.X] = true
        }
    }

    return m
}

func clear() {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }

    cmd.Stdout = os.Stdout
    err := cmd.Run()
    if err != nil {
        panic(err)
    }
}
