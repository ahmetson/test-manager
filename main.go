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

	err = socket.Bind("tcp://*:6000")
	if err != nil {
		println("bind error")
		panic(err)
	}

	time.Sleep(time.Second * N)

	err = socket.Close()
	if err != nil {
		println("close error")
		panic(err)
	}

	println("Server exited successfully")
}
