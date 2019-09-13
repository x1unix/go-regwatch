// +build !windows

package regwatch

import (
	"fmt"
	"runtime"
)

// Key represents registry key
type Key uintptr

// NewWatcher creates a new watcher instance
func NewWatcher(hMainKey Key, regPath string, timeout Timeout) (Watcher, error) {
	return nil, fmt.Printf("this package is not supported on platform %q", runtime.GOOS)
}
