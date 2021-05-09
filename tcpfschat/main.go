package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    str := os.Args[1]
    fmt.Println(filepath.Base(str))
}
