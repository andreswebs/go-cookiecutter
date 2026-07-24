// Package terr implements typed, coded errors: every failure that can reach a
// consumer carries a stable machine code, a process exit code, and an optional
// remediation hint, discovered at the exit boundary via errors.As against the
// Coded interface.
//
// Sentinels are immutable package-level values created with New, which also
// registers them so the tool's error inventory is enumerable via All (this is
// what the schema command reports). Per-invocation context attaches by copy:
// Wrap and WithDetails return modified copies and never mutate the receiver,
// so sentinels are safe to share across goroutines by construction.
package terr

import "fmt"

// Coded is an error that carries a stable machine code, a process exit code,
// and a user-facing remediation hint. The exit boundary resolves it with
// errors.As to render the error envelope and choose the exit code.
type Coded interface {
	error
	Code() string
	ExitCode() int
	Hint() string
}

// Detailed is an optional interface for errors that carry render-ready
// structured details, surfaced in the error envelope's "details" field.
type Detailed interface {
	ErrorDetails() any
}

// E is the concrete coded error. Declare sentinels with New, attach a cause
// with Wrap, and attach structured payloads with WithDetails.
type E struct {
	code, msg, hint string
	exit            int
	cause           error
	details         any
}

var (
	_ Coded    = (*E)(nil)
	_ Detailed = (*E)(nil)
)

var registry []*E

// New creates a sentinel E and registers it for enumeration via All. It
// panics when code is already registered: duplicate registration is an
// init-time programmer error, and crashing at startup is the correct outcome.
func New(code string, exit int, hint, msg string) *E {
	for _, r := range registry {
		if r.code == code {
			panic(fmt.Sprintf("terr: duplicate error code %q", code))
		}
	}
	e := &E{code: code, exit: exit, hint: hint, msg: msg}
	registry = append(registry, e)
	return e
}

// Newf creates an E without registering it, formatting the message from
// format and args. Use it for one-off per-invocation errors whose class does
// not belong in the enumerable inventory.
func Newf(code string, exit int, hint, format string, args ...any) *E {
	return &E{code: code, exit: exit, hint: hint, msg: fmt.Sprintf(format, args...)}
}

// All returns a copy of every error registered via New, in registration
// order. It backs the schema command's error inventory and the exit-code
// conformance test.
func All() []*E {
	cp := make([]*E, len(registry))
	copy(cp, registry)
	return cp
}

// Error returns the message; when a cause is attached it is appended as
// "message: cause".
func (e *E) Error() string {
	if e.cause != nil {
		return e.msg + ": " + e.cause.Error()
	}
	return e.msg
}

// Code returns the stable machine code.
func (e *E) Code() string { return e.code }

// ExitCode returns the process exit code associated with the error.
func (e *E) ExitCode() int { return e.exit }

// Hint returns the user-facing remediation hint, or "" when there is none.
func (e *E) Hint() string { return e.hint }

// Unwrap returns the attached cause, or nil, so errors.Is and errors.As
// traverse the chain.
func (e *E) Unwrap() error { return e.cause }

// Is reports whether target is an E with the same code, so copies produced by
// Wrap and WithDetails still match their sentinel under errors.Is.
func (e *E) Is(target error) bool {
	t, ok := target.(*E)
	return ok && t.code == e.code
}

// Wrap returns a copy of e with cause attached as its underlying error. The
// receiver is left unchanged so sentinels stay immutable.
func (e *E) Wrap(cause error) *E {
	cp := *e
	cp.cause = cause
	return &cp
}

// WithDetails returns a copy of e carrying the given structured details. The
// receiver is left unchanged so sentinels stay immutable.
func (e *E) WithDetails(details any) *E {
	cp := *e
	cp.details = details
	return &cp
}

// ErrorDetails returns the structured details attached with WithDetails, or
// nil.
func (e *E) ErrorDetails() any { return e.details }
