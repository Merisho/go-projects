package server

import (
    "net"
    "strconv"
)

func ServeTCP(port uint16) (*Server, error) {
    listener, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(port), 10))
    if err != nil {
        return &Server{}, err
    }

    return NewServer(listener), nil
}
