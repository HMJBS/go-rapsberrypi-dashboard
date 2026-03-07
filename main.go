// Package main provides the dashboard executable.
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Hello, Raspberry Pi 1!")
	fmt.Printf("Running on: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
