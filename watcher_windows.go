package regwatch

/*
#cgo CFLAGS: -I.
#define DEBUG 0
#include "watch.h"
*/
import "C"

import (
	"syscall"
	"unsafe"
)

const (
    // HKeyClassesRoot represents HKEY_CLASSES_ROOT hive
	HKeyClassesRoot     = Key(syscall.HKEY_CLASSES_ROOT)
	
	// HKeyCurrentUser represents HKEY_CURRENT_USER hive
	HKeyCurrentUser     = Key(syscall.HKEY_CURRENT_USER)
	
	// HKeyLocalMachine represents HKEY_LOCAL_MACHINE hive
	HKeyLocalMachine    = Key(syscall.HKEY_LOCAL_MACHINE)
	
	// HKeyUsers represents HKEY_USERS hive
	HKeyUsers           = Key(syscall.HKEY_USERS)
	
	// HKeyCurrentConfig represents HKEY_CURRENT_CONFIG hive
	HKeyCurrentConfig   = Key(syscall.HKEY_CURRENT_CONFIG)
	
	// HKeyPerformanceData represents HKEY_PERFORMANCE_DATA hive
	HKeyPerformanceData = Key(syscall.HKEY_PERFORMANCE_DATA)
	
	// Infinity is infinite timeout
	Infinity = 0xFFFFFFFF
)

// Key represents registry key
type Key syscall.Handle

type watcherProxy struct {
	w C.Watcher
	timeout Timeout
}

func (w *watcherProxy) Destroy() error {
	errno := C.watcher_close(&w.w)
	if errno > 0 {
		return syscall.Errno(errno)
	}

	return nil
}

func (w *watcherProxy) Await() (bool, error) {
	var changed C.uchar
	errno := C.watcher_await(&w.w, C.long(w.timeout), &changed)
	if errno > 0 {
		return false, syscall.Errno(errno)
	}

	return changed > 0, nil
}


// NewWatcher creates a new watcher instance
func NewWatcher(hMainKey Key, regPath string, timeout Timeout) (Watcher, error) {
	var w C.Watcher
	cstr := C.CString(regPath)
	errno := C.watcher_create(unsafe.Pointer(hMainKey), cstr, &w);
	if errno > 0 {
		return nil, syscall.Errno(errno)
	}

	return &watcherProxy{w: w, timeout: timeout}, nil
}