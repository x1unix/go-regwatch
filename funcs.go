package main

import (
	"fmt"
	"os"
)

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
