package main

import (
	"os"
	"time"
)

func main() {
	println("hello world!")

	time.Sleep(time.Millisecond * 100)

	// it exits with an error
	os.Exit(1)
}
