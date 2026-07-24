package command

import "strings"

// classifyUsage is the single sanctioned string-matching classifier: the CLI
// framework's parse errors carry no types, so this is the one place errors
// are classified by message instead of identity. It returns err wrapped in
// ErrUsage when the message matches a known framework parse failure, nil
// otherwise. The golden scenarios that drive real parse failures through Run
// fail when the framework's wording changes, keeping this list honest.
func classifyUsage(err error) error {
	msg := err.Error()
	for _, prefix := range []string{
		"flag provided but not defined",
		"flag needs an argument",
		"invalid value",
		"could not parse",
		"Required flag",
		"Required flags",
	} {
		if strings.HasPrefix(msg, prefix) {
			return ErrUsage.Wrap(err)
		}
	}
	if strings.HasPrefix(msg, "option ") && strings.Contains(msg, "cannot be set along with option") {
		return ErrUsage.Wrap(err)
	}
	return nil
}
