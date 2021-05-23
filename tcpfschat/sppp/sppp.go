package sppp

type MessageType byte

const (
    TextType = iota + 1
    StreamType
    EndType
    ErrorType
    TimeoutType
    maxMessageTypeIota
)

func Uint64ToBytes(n uint64) [8]byte {
    if n == 0 {
        return [8]byte{}
    }

    var res [8]byte
    i := 0
    for n > 0 {
        res[i] = byte(n % 256)
        n /= 256
        i++
    }

    return res
}

func BytesToUint64(b [8]byte) uint64 {
    var n uint64
    var p uint64 = 1
    for i := 0; i < 8; i++ {
        n += uint64(b[i]) * p
        p *= 256
    }

    return n
}
