package errors

import "fmt"

// -----------------------------------------------------------------------------
// Public Errors
// -----------------------------------------------------------------------------

var (
	// ErrUnimplemented indicates the functionality requested is not yet implemented
	ErrUnimplemented = fmt.Errorf("unimplemented")

	// ErrMissing indicates a requested resource is not present
	ErrMissing = fmt.Errorf("missing")

	// ErrExists indicates that a resource requested for creation already exists
	ErrExists = fmt.Errorf("exists")
)
