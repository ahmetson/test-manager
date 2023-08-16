package main

import (
	"os"
	"time"
)

func main() {
	println("hello world!")

	time.Sleep(time.Second)

	// it exits with an error
	os.Exit(1)
}
