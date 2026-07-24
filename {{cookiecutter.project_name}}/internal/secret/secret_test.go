package secret_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/{{ cookiecutter.author_id }}/{{ cookiecutter.project_name }}/internal/secret"
)

const raw = "hunter2-super-secret"

func TestRedactsOnAllLeakSurfaces(t *testing.T) {
	v := secret.Value(raw)

	for _, verb := range []string{"%v", "%s", "%q", "%#v"} {
		got := fmt.Sprintf(verb, v)
		if strings.Contains(got, raw) {
			t.Errorf("fmt %s leaked the value: %s", verb, got)
		}
		if !strings.Contains(got, "REDACTED") {
			t.Errorf("fmt %s did not redact: %s", verb, got)
		}
	}

	data, err := json.Marshal(struct {
		Token secret.Value `json:"token"`
	}{Token: v})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if strings.Contains(string(data), raw) {
		t.Errorf("json.Marshal leaked the value: %s", data)
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Info("connecting", "token", v)
	if strings.Contains(buf.String(), raw) {
		t.Errorf("slog leaked the value: %s", buf.String())
	}
}

func TestRevealReturnsTheRawValue(t *testing.T) {
	if got := secret.Value(raw).Reveal(); got != raw {
		t.Errorf("Reveal() = %q, want %q", got, raw)
	}
}
