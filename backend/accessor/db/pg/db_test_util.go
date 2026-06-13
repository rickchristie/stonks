package pg

import (
	"context"
	"time"
)

// TimeoutOverrideForTesting should only be used in tests. Long-running
// application transactions are usually a design smell.
func TimeoutOverrideForTesting(parentCtx context.Context, dur time.Duration) context.Context {
	return context.WithValue(parentCtx, timeoutKey, dur)
}
