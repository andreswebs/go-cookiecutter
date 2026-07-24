package command

// ExitCodes declares, as data, every exit code this tool can produce. The
// conformance tests assert that real exits stay inside this registry and
// that every declared code is exercised by at least one test scenario, so
// the declaration cannot drift from behavior.
//
//	0  success, clean
//	64 usage error (EX_USAGE): the CLI surface was misused
//	65 data error (EX_DATAERR): input payload rejected
//	70 internal error (EX_SOFTWARE): a bug, including any unclassified error
//
// Result classes (exit 1, and optional sub-codes 2-63) are not declared
// because the placeholder commands produce none; add them here when a
// command gains a result that demands action.
var ExitCodes = []int{0, 64, 65, 70}
