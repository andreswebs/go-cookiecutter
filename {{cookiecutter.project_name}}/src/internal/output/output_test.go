package output_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/output"
	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/terr"
)

func TestEmitJSONAppendsOneNewline(t *testing.T) {
	var buf bytes.Buffer
	env := output.GreetEnvelope{Head: output.OKHead(), Greeting: "hello, world"}

	if err := output.EmitJSON(&buf, env); err != nil {
		t.Fatalf("EmitJSON: %v", err)
	}
	got := buf.String()
	if !strings.HasSuffix(got, "\n") || strings.Count(got, "\n") != 1 {
		t.Errorf("output is not one newline-terminated object: %q", got)
	}
	if !strings.HasPrefix(got, `{"schema_version":1,"ok":true`) {
		t.Errorf("envelope does not open with the head: %q", got)
	}
}

func TestEmitErrorRendersCodedFields(t *testing.T) {
	sentinel := terr.Newf("bad_input", 65, "fix the input", "input rejected")
	var buf bytes.Buffer

	output.EmitError(&buf, sentinel.WithDetails(map[string]string{"field": "name"}))

	var env struct {
		SchemaVersion int  `json:"schema_version"`
		OK            bool `json:"ok"`
		Error         struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Hint    string         `json:"hint"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("error envelope is not valid JSON: %v", err)
	}
	if env.OK {
		t.Error("error envelope has ok=true")
	}
	if env.Error.Code != "bad_input" || env.Error.Hint != "fix the input" {
		t.Errorf("code/hint not resolved: %+v", env.Error)
	}
	if env.Error.Message != "input rejected" {
		t.Errorf("message = %q", env.Error.Message)
	}
	if env.Error.Details["field"] != "name" {
		t.Errorf("details not rendered: %+v", env.Error.Details)
	}
}

func TestEmitErrorUnclassifiedIsInternal(t *testing.T) {
	var buf bytes.Buffer
	output.EmitError(&buf, errors.New("boom"))

	if !strings.Contains(buf.String(), `"code":"internal_error"`) {
		t.Errorf("unclassified error not rendered as internal_error: %s", buf.String())
	}
}

func TestExitCodeFor(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"coded error carries its own", terr.Newf("x", 65, "", "x"), 65},
		{"unclassified is internal", errors.New("boom"), 70},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := output.ExitCodeFor(tt.err); got != tt.want {
				t.Errorf("ExitCodeFor() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEmitWarningIsNDJSON(t *testing.T) {
	var buf bytes.Buffer
	output.EmitWarning(&buf, "deprecated_flag", "flag renamed", "use --new-name", nil)
	output.EmitWarning(&buf, "slow_path", "fallback engaged", "", nil)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 NDJSON lines, got %d: %q", len(lines), buf.String())
	}
	for _, line := range lines {
		var env struct {
			Level string `json:"level"`
		}
		if err := json.Unmarshal([]byte(line), &env); err != nil {
			t.Fatalf("warning line is not valid JSON: %v", err)
		}
		if env.Level != "warning" {
			t.Errorf("level = %q, want warning", env.Level)
		}
	}
}

func TestSchemaEnvelopeCoalescesCollections(t *testing.T) {
	data, err := json.Marshal(output.SchemaEnvelope{Head: output.OKHead()})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	got := string(data)
	for _, want := range []string{`"commands":[]`, `"exit_codes":[]`, `"errors":[]`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "null") {
		t.Errorf("collection serialized as null: %s", got)
	}
}
