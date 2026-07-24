// Package output defines the machine-readable stream contract: the result
// envelope head shared by every command, the JSON error envelope rendered on
// stderr, and NDJSON warning lines. It is framework-free and never touches
// the real process streams; callers inject writers.
package output

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"
)

// SchemaVersion is the version of the output contract, bumped on breaking
// shape changes to any envelope.
const SchemaVersion = 1

// Head opens every result envelope: embed it as the first field of each
// command's envelope struct.
type Head struct {
	SchemaVersion int  `json:"schema_version"`
	OK            bool `json:"ok"`
}

// OKHead returns the envelope head for a successful result.
func OKHead() Head {
	return Head{SchemaVersion: SchemaVersion, OK: true}
}

// EmitJSON marshals v and writes it to w followed by a newline: one
// newline-terminated JSON object per invocation.
func EmitJSON(w io.Writer, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

type errorEnvelope struct {
	SchemaVersion int         `json:"schema_version"`
	OK            bool        `json:"ok"`
	Error         errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
	Details any    `json:"details,omitempty"`
}

// EmitError writes the JSON error envelope for err to w, which is the stderr
// sink. The code and hint come from the typed error in err's chain; an
// unclassified error renders as internal_error so missing classifications
// stay visible. The write is best-effort: a failure to report an error on
// stderr is unrecoverable, so it is not escalated over the error it
// describes.
func EmitError(w io.Writer, err error) {
	env := errorEnvelope{
		SchemaVersion: SchemaVersion,
		Error:         errorDetail{Code: "internal_error", Message: err.Error()},
	}

	var coded terr.Coded
	if errors.As(err, &coded) {
		env.Error.Code = coded.Code()
		env.Error.Hint = coded.Hint()
	}
	var detailed terr.Detailed
	if errors.As(err, &detailed) {
		env.Error.Details = detailed.ErrorDetails()
	}

	data, merr := json.Marshal(env)
	if merr != nil {
		// Unmarshalable details: degrade to the envelope without them,
		// which cannot fail to marshal.
		env.Error.Details = nil
		data, _ = json.Marshal(env)
	}
	_, _ = w.Write(append(data, '\n'))
}

// ExitCodeFor maps err to its process exit code. Typed errors carry their
// own; anything unclassified is an internal error, exit 70 (EX_SOFTWARE), so
// an unclassified failure is loud instead of silently mislabeled.
func ExitCodeFor(err error) int {
	var coded terr.Coded
	if errors.As(err, &coded) {
		return coded.ExitCode()
	}
	return 70
}

// warningEnvelope is a non-fatal, machine-readable advisory written to
// stderr. It carries level "warning" instead of an ok field, so a consumer
// can tell it apart from the error envelope unambiguously, and it never
// changes the exit code.
type warningEnvelope struct {
	SchemaVersion int    `json:"schema_version"`
	Level         string `json:"level"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	Hint          string `json:"hint,omitempty"`
	Details       any    `json:"details,omitempty"`
}

// EmitWarning writes one JSON warning line to w, which is the stderr sink.
// Each call appends exactly one newline-terminated object, so successive
// warnings form a valid NDJSON stream. Libraries never call this directly:
// they raise advisories through an injected callback and stay stream-blind.
// The write is best-effort, like EmitError.
func EmitWarning(w io.Writer, code, message, hint string, details any) {
	env := warningEnvelope{
		SchemaVersion: SchemaVersion,
		Level:         "warning",
		Code:          code,
		Message:       message,
		Hint:          hint,
		Details:       details,
	}
	data, err := json.Marshal(env)
	if err != nil {
		return
	}
	_, _ = w.Write(append(data, '\n'))
}
