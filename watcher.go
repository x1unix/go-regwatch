// Package regwatch provides a tiny wrapper that allows to track registry key changes
// in Windows operating system.
//
// See "example_test.go" for usage example.
package regwatch

// Timeout represents event timeout in milliseconds
type Timeout uint32


// Watcher tracks registry changes
type Watcher interface {
	// Destroy stops
	Destroy() error

	// Await waits until key change event occurs
	Await() (bool, error)
}