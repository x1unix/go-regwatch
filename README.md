# go-regwatch

Package regwatch provides a tiny wrapper that allows to track registry key changes
in Windows operating system.

This library wraps `RegNotifyChangeKeyValue` Windows syscall.

## Usage

See `example_test.go`