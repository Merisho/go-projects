package csdt // concurrency safe data type

import "sync"

func NewValue() *Value {
    return &Value{}
}

type Value struct {
    val interface{}
    m sync.RWMutex
}

func (v *Value) Value() interface{} {
    v.m.RLock()
    defer v.m.RUnlock()

    return v.val
}

func (v *Value) SetValue(val interface{}) {
    v.m.Lock()
    defer v.m.Unlock()

    v.val = val
}
