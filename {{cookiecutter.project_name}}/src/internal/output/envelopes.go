package output

import "encoding/json"

// GreetEnvelope carries the result of the placeholder greet command. Replace
// it, and add one envelope struct per command, each opening with Head.
type GreetEnvelope struct {
	Head
	Greeting string `json:"greeting"`
}

// SchemaEnvelope carries the tool's runtime self-description: the command
// surface, the declared exit codes, and the error inventory. It is a
// projection of the declarations that drive behavior, never hand-maintained.
type SchemaEnvelope struct {
	Head
	Tool      string          `json:"tool"`
	Version   string          `json:"version"`
	Commands  []SchemaCommand `json:"commands"`
	ExitCodes []int           `json:"exit_codes"`
	Errors    []SchemaError   `json:"errors"`
}

// SchemaCommand describes one command: its name, usage line, argument
// synopsis, and flags.
type SchemaCommand struct {
	Name  string       `json:"name"`
	Usage string       `json:"usage"`
	Args  string       `json:"args,omitempty"`
	Flags []SchemaFlag `json:"flags"`
}

// SchemaFlag describes one flag of a command.
type SchemaFlag struct {
	Name  string `json:"name"`
	Usage string `json:"usage,omitempty"`
}

// SchemaError describes one entry of the error inventory: a stable machine
// code, the exit code it maps to, and its remediation hint.
type SchemaError struct {
	Code     string `json:"code"`
	ExitCode int    `json:"exit_code"`
	Hint     string `json:"hint,omitempty"`
}

// MarshalJSON emits the envelope with non-nil collections: absent lists are
// [], enforced here rather than by call-site discipline.
func (e SchemaEnvelope) MarshalJSON() ([]byte, error) {
	type alias SchemaEnvelope
	a := alias(e)
	if a.Commands == nil {
		a.Commands = []SchemaCommand{}
	}
	if a.ExitCodes == nil {
		a.ExitCodes = []int{}
	}
	if a.Errors == nil {
		a.Errors = []SchemaError{}
	}
	return json.Marshal(a)
}

// MarshalJSON emits the command description with a non-nil flags list.
func (c SchemaCommand) MarshalJSON() ([]byte, error) {
	type alias SchemaCommand
	a := alias(c)
	if a.Flags == nil {
		a.Flags = []SchemaFlag{}
	}
	return json.Marshal(a)
}
