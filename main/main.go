package main

import (
	"fmt"
	"rec"
	"runtime"
)

func main() {
	rec.Printf("self", "Go version: %s\n", runtime.Version())
	rec.Printf("self", "OS: %s\n", runtime.GOOS)
	rec.Printf("self", "Architecture: %s\n", runtime.GOARCH)
	rec.Printf("self", "Compiler: %s\n", runtime.Compiler)
	var a any
	test(nil)
}

func test(str fmt.Stringer) {

}
