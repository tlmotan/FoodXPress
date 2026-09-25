package domain

import "errors"

// Domain-level sentinels. Repositories wrap these with %w so callers can
// classify a failure (404 vs 500) without importing the repository package.
var (
	ErrDriverNotFound = errors.New("driver not found")
)
