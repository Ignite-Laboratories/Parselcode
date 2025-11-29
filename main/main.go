package main

import (
	"rec"
	"runtime"
)

func main() {
	rec.Printf("self", "Go version: %s\n", runtime.Version())
	rec.Printf("self", "OS: %s\n", runtime.GOOS)
	rec.Printf("self", "Architecture: %s\n", runtime.GOARCH)
	rec.Printf("self", "Compiler: %s\n", runtime.Compiler)

}

func MoveBy42[T any](vec T) (out T, err error) {

}
