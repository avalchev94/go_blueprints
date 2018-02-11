package trace

import (
	"fmt"
	"io"
	"time"
)

// Tracer is the interface that describes an object capable of tracing
// events throught code.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (tr *tracer) Trace(a ...interface{}) {
	t := time.Now()

	tr.out.Write([]byte(t.String() + ": "))
	tr.out.Write([]byte(fmt.Sprint(a...)))
	tr.out.Write([]byte("\n"))
}

// New creates tracer object
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}
