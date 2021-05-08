package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "net"
    "strconv"
)

func Serve(port uint16) (*Server, error) {
    listener, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(port), 10))
    if err != nil {
        return nil, err
    }

    sp := sppp.NewSPPPListener(listener)

    return newServer(sp), nil
}
