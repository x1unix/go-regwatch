package regwatch

import (
	"fmt"
	"sync"
	"os"
	"os/signal"
	"context"
	"runtime"
	"testing"
	"time"
)

const keyPath = `SOFTWARE\SAAZOD\ManagedPosix`

func TestExample(t *testing.T) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx, cancelFn := cancelHandler()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	w, err := NewWatcher(HKeyLocalMachine, keyPath, 1000)
	must(err)

	go func() {
		<-time.After(5 * time.Second)
		fmt.Printf("INFO:\tstopping...\n")
		cancelFn()
		os.Exit(1)
	}()

	updates := make(chan string)
	go func() {
		fmt.Printf("INFO:\tconsumer: start...\n")
		for {
			kp, ok := <-updates
			if !ok {
				fmt.Printf("INFO:\tconsumer: shutdown...\n")
				return
			}

			fmt.Printf("INFO:\t consumer: Received update msg - '%s'\n", kp)
		}
	}()

	go func() {
		fmt.Printf("INFO:\tLooking for changes in '%s'...\n", keyPath)
		defer func() {
			if err := w.Destroy(); err != nil {
				fmt.Printf("ERROR:\tWatcher.Destroy - %s\n", err)
			}

			close(updates)
			wg.Done();
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("INFO:\tGot shutdown signal...\n")
				return
			default:
			}

			changed, err := w.Await()
			if err != nil {
				fmt.Printf("ERROR:\tWatcher.Await - %s\n", err)
				return
			}

			if !changed {
				continue
			}

			fmt.Printf("INFO:\t'%s' changed\n", keyPath)
			updates <- keyPath
		}
	}()


	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Goodbye")
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func try(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
	}
}