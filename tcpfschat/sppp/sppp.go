package sppp

type MessageType byte

const (
    TextType = iota
    FileType
)

func Int64ToBytes(n int64) [8]byte {
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

func BytesToInt64(b [8]byte) int64 {
    var n int64
    var p int64 = 1
    for i := 0; i < 8; i++ {
        n += int64(b[i]) * p
        p *= 256
    }

    return n
}
