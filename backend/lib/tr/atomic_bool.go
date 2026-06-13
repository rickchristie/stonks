package tr

import "sync/atomic"

type AtomicBool struct {
	v atomic.Bool
}

func (b *AtomicBool) Store(v bool) {
	b.v.Store(v)
}

func (b *AtomicBool) Load() bool {
	return b.v.Load()
}
