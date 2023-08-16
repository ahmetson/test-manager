package main

import (
	"github.com/pebbe/zmq4"
	"time"
)

const N = 2

// For N seconds, it imitates the SDS handler as busy.
// For testing dep_manager.Running() command.
func main() {

	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		panic(err)
	}

	// port = 8081
	err = socket.Bind("tpc://*:8081")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * N)

	err = socket.Close()
	if err != nil {
		panic(err)
	}

	println("Server exited successfully")
}
