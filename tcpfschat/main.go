package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        ln, err := net.Listen("tcp", "localhost:1337")
        if err != nil {
            panic(err)
        }

        wg.Done()
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }

        for {
            time.Sleep(1 * time.Second)
            _, err := conn.Write([]byte("Hello"))
            if err != nil {
                fmt.Println(err.(*net.OpError).Error())
                break
            }

            fmt.Println("Sent Hello")
        }

        _, err = conn.Read([]byte{})
        fmt.Println(err, err.(*net.OpError).Temporary())
    }()

    wg.Wait()

    conn, err := net.Dial("tcp", "localhost:1337")
    if err != nil {
        panic(err)
    }

    k := 0
    b := make([]byte, 128)
    for {
        if k == 3 {
            conn.Close()
            break
        }
        n, err := conn.Read(b)
        if err != nil {
            panic(err)
        }

        fmt.Println(string(b[:n]))
        k++
    }

    select {}
}
