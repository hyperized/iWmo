package iwmo

// Re-exports of private symbols for use by the external (iwmo_test) test
// package. Compile-time visibility is controlled by the _test.go suffix:
// these definitions are only visible while running `go test`, so the public
// API is unaffected.

// Internal validation helpers exposed for external tests.
var (
	ValidateHeader   = validateHeader
	ValidateGeslacht = validateGeslacht
)

// Fixtures used by tests in the external (iwmo_test) package. Internal tests
// can keep calling the lowercase helpers in fixtures_test.go directly.
var (
	ValidHeaderFixture = validHeader
	ValidWMO301        = validWMO301
	ValidWMO302        = validWMO302
	ValidWMO303        = validWMO303
	ValidWMO304        = validWMO304
	ValidWMO305        = validWMO305
	ValidWMO315        = validWMO315
)

// ValidBSN is a BSN known to pass the elfproef. Re-exported for external tests.
const ValidBSN = validBSN
