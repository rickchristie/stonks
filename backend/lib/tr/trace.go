package tr

import "time"

// Trace carries request context through service and accessor layers. Keep it
// small and scrub sensitive request fields before attaching them.
type Trace struct {
	TraceId   string
	Start     time.Time
	Path      string
	LogDetail AtomicBool
	Request   any
}

func New(traceId string) *Trace {
	return &Trace{
		TraceId: traceId,
		Start:   time.Now(),
	}
}
