package mend

import (
	"log"
	"strconv"

	"app-template/lib/tr"
)

type Logger interface {
	FatalErr(trace *tr.Trace, err error) Log
	ErrorErr(trace *tr.Trace, err error) Log
	WarnErr(trace *tr.Trace, err error) Log
	InfoNoTrace() Log
	ErrorNoTrace() Log
	WarnNoTrace() Log
}

type Log interface {
	Str(key string, val string) Log
	Int(key string, val int) Log
	Error(err error) Log
	Msg(msg string)
}

type StdLogger struct {
	name string
}

func NewZerologLogger(name string) *StdLogger {
	return &StdLogger{name: name}
}

func (l *StdLogger) FatalErr(trace *tr.Trace, err error) Log {
	return l.event("fatal", trace).Error(err)
}

func (l *StdLogger) ErrorErr(trace *tr.Trace, err error) Log {
	return l.event("error", trace).Error(err)
}

func (l *StdLogger) WarnErr(trace *tr.Trace, err error) Log {
	return l.event("warn", trace).Error(err)
}

func (l *StdLogger) InfoNoTrace() Log {
	return l.event("info", nil)
}

func (l *StdLogger) ErrorNoTrace() Log {
	return l.event("error", nil)
}

func (l *StdLogger) WarnNoTrace() Log {
	return l.event("warn", nil)
}

func (l *StdLogger) event(level string, trace *tr.Trace) *StdLog {
	e := &StdLog{
		level: level,
		name:  l.name,
		vals:  map[string]string{},
	}
	if trace != nil {
		e.vals["_tid"] = trace.TraceId
	}
	return e
}

type StdLog struct {
	level string
	name  string
	err   error
	vals  map[string]string
}

func (l *StdLog) Str(key string, val string) Log {
	l.vals[key] = val
	return l
}

func (l *StdLog) Int(key string, val int) Log {
	l.vals[key] = strconv.Itoa(val)
	return l
}

func (l *StdLog) Error(err error) Log {
	l.err = err
	return l
}

func (l *StdLog) Msg(msg string) {
	if l.err != nil {
		log.Printf("level=%s logger=%s msg=%q err=%q vals=%v", l.level, l.name, msg, l.err.Error(), l.vals)
		return
	}
	log.Printf("level=%s logger=%s msg=%q vals=%v", l.level, l.name, msg, l.vals)
}
