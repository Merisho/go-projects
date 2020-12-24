package chaterrors

import (
    "net"
)

func IsTemporary(err error) bool {
    if err == nil {
        return true
    }

    if e, ok := err.(net.Error); ok {
        return e.Temporary()
    }

    return false
}
