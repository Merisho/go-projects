package main

import (
    "bytes"
    "fmt"
)

func main() {
    b := []byte{1,1,0,1,1,0}
    bs := bytes.Split(b, []byte{0})

    fmt.Println(len(bs), bs)
}
