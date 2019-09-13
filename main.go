package main

/*
#cgo CFLAGS: -I.
#include "watch.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"syscall"
	"os"
	"os/signal"
	"context"
	"runtime"
	"time"
)

const keyPath = `SOFTWARE\SAAZOD\ManagedPosix`

type listener func(regPath string, event uint32) bool

var listeners = map[string]listener{}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx, _ := cancelHandler()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(time.Millisecond)
		fmt.Printf("INFO:\tLooking for changes in '%s'...\n", keyPath)
		defer wg.Done();
		select {
		case <-ctx.Done():
			fmt.Printf("INFO:\tGot shutdown signal...\n")
		default:
		}
		err := RegListen(C.HKEY_LOCAL_MACHINE, keyPath, func(kp string, event uint32) bool {
			fmt.Printf("INFO:\t'%s' changed (%d)\n", kp, event)
			select {
			case <-ctx.Done():
				fmt.Printf("DEBUG:\t context is dead\n")
				return false
			default:
				return true
			}
		})
		if err != nil {
			fmt.Printf("ERR:\t%s", err)
		}
	}(wg)

	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Goodbye")
}

func RegListen(hKey C.HKEY, regPath string, cb listener) error {
	cstr := C.CString(regPath)
	listeners[regPath] = cb
	result := C.reg_listen(hKey, cstr)
	if result == C.ERROR_SUCCESS {
		return nil
	}

	return syscall.Errno(result)
}

//export reg_global_listener
func reg_global_listener(cstr *C.char, dwEvent uint32) bool {
	keyPath := C.GoString(cstr)
	handler, ok := listeners[keyPath]
	if !ok {
		fmt.Printf("ERROR:\tno listener for '%s'\n", keyPath)
		return false;
	}

	result := handler(keyPath, dwEvent)
	return result
}


// cancelHandler returns cancellation context and function for graceful shutdown
func cancelHandler() (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancelFn()
		signal.Stop(signals)
	}()

	return ctx, cancelFn
}
