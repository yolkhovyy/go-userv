//go:build test

package main

import (
	"os"
	"testing"
)

func TestRunMain(t *testing.T) {
	os.Args = []string{os.Args[0]}
	run()
}
