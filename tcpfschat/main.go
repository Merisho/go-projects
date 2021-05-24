package main

import (
    "bytes"
    "os"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(3)

    go func() {
        cf("a.txt")
        wg.Done()
    }()

    go func() {
        cf("b.txt")
        wg.Done()
    }()

    go func() {
        cf("c.txt")
        wg.Done()
    }()

    wg.Wait()
}

func cf(name string) {
    f, err := os.OpenFile("./" + name, os.O_CREATE, 0777)
    if err != nil {
        panic(err)
    }

    for i := 0; i < 1024; i++ {
        b := bytes.Repeat([]byte(name),  1024 * 1024)
        _, err := f.Write(b)
        if err != nil {
            panic(err)
        }
    }

    err = f.Close()
    if err != nil {
        panic(err)
    }
}
