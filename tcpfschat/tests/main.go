package main

import "fmt"

func main() {
    c := make(chan int, 10)
    for i := 0; i < 10; i++ {
        c <- i
        fmt.Println(len(c))
    }

    for len(c) > 0 {
        fmt.Println("From channel:", <-c)
    }
}
