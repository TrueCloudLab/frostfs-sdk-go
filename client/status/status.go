package apistatus

// Status defines a variety of FrostFS API status returns.
//
// All statuses are split into two disjoint subsets: successful and failed, and:
//   - statuses that implement the build-in error interface are considered failed statuses;
//   - all other value types are considered successes (nil is a default success).
//
// In Go code type of success can be determined by a type switch, failure - by a switch with errors.As calls.
// Nil should be considered as a success, and default switch section - as an unrecognized Status.
//
// To convert statuses into errors and vice versa, use functions ErrToStatus and ErrFromStatus, respectively.
// ErrFromStatus function returns nil for successful statuses. However, to simplify the check of statuses for success,
// IsSuccessful function should be used (try to avoid nil comparison).
// It should be noted that using direct typecasting is not a compatible approach.
//
// To transport statuses using the FrostFS API V2 protocol, see StatusV2 interface and FromStatusV2 and ToStatusV2 functions.
type Status interface{}

// ErrFromStatus converts Status instance to error if it is failed. Returns nil on successful Status.
//
// Note: direct assignment may not be compatibility-safe.
func ErrFromStatus(st Status) error {
	if err, ok := st.(error); ok {
		return err
	}

	return nil
}

// ErrToStatus converts the error instance to Status instance.
//
// Note: direct assignment may not be compatibility-safe.
func ErrToStatus(err error) Status {
	return err
}

// IsSuccessful checks if status is successful.
//
// Note: direct cast may not be compatibility-safe.
func IsSuccessful(st Status) bool {
	_, ok := st.(error)
	return !ok
}
