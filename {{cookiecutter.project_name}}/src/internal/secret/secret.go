// Package secret keeps credential values structurally out of logs and
// serialized output. Any value that is a credential lives in the Value type;
// the raw string is reachable only through the explicit, greppable Reveal
// method.
package secret

import "log/slog"

// redactedMarker is the placeholder rendered in place of a Value across every
// logging and serialization path.
const redactedMarker = "REDACTED"

// Value is a string that redacts itself on all four leak surfaces: slog
// (LogValue), fmt's %v/%s/%q verbs (String), fmt's %#v verb (GoString), and
// encoding/json (MarshalJSON). Accidentally logging or marshaling a struct
// that embeds it can never leak the credential; only Reveal returns the raw
// string.
type Value string

// LogValue redacts the value when logged through slog.
func (Value) LogValue() slog.Value { return slog.StringValue(redactedMarker) }

// String redacts the value under fmt's %v, %s, and %q verbs.
func (Value) String() string { return redactedMarker }

// GoString redacts the value under fmt's %#v verb, which would otherwise
// print the underlying value of a defined string type.
func (Value) GoString() string { return redactedMarker }

// MarshalJSON redacts the value when a containing struct is JSON-encoded,
// including by the slog JSON handler.
func (Value) MarshalJSON() ([]byte, error) { return []byte(`"` + redactedMarker + `"`), nil }

// Reveal returns the underlying credential. It is the only way to obtain the
// raw string and must be called only where the value is actually required.
func (v Value) Reveal() string { return string(v) }
