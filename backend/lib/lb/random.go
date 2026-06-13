package lb

import (
	"crypto/rand"
	"encoding/binary"
	"sync/atomic"
	"time"
)

type randomizedEndpoints[T any] struct {
	endpoints []T
	counter   uint64
}

func Randomized[T any](endpoints []T) Selector[T] {
	if len(endpoints) == 0 {
		panic("endpoint list cannot be empty")
	}

	endpointsCopy := make([]T, len(endpoints))
	copy(endpointsCopy, endpoints)

	var initialSeed uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &initialSeed); err != nil {
		initialSeed = uint64(time.Now().UnixNano())
	}

	return &randomizedEndpoints[T]{
		endpoints: endpointsCopy,
		counter:   initialSeed,
	}
}

func (r *randomizedEndpoints[T]) Get() T {
	if len(r.endpoints) == 1 {
		return r.endpoints[0]
	}

	current := atomic.AddUint64(&r.counter, 1)
	seed := current*1664525 + 1013904223
	seed ^= seed >> 13
	seed ^= seed << 17
	seed ^= seed >> 5

	return r.endpoints[seed%uint64(len(r.endpoints))]
}

func (r *randomizedEndpoints[T]) Count() int {
	return len(r.endpoints)
}
