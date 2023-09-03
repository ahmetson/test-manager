package main

import (
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/common-lib/message"
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

	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)

	currentTime := time.Now()
	endTime := currentTime.Add(time.Second * N)

	for {
		polled, err := poller.Poll(time.Millisecond)
		if err != nil {
			panic(err)
		}

		currentTime = currentTime.Add(time.Millisecond)
		if currentTime.UnixMilli() >= endTime.UnixMilli() {
			break
		}

		if len(polled) == 0 {
			continue
		}

		_, err = socket.RecvMessage(0)
		if err != nil {
			panic(err)
		}

		reply := message.Reply{
			Status:     message.OK,
			Parameters: key_value.Empty(),
			Message:    "",
		}
		replyStr, err := reply.String()
		if err != nil {
			panic(err)
		}
		_, err = socket.SendMessage(replyStr)
		if err != nil {
			panic(err)
		}
		break
	}

	err = socket.Close()
	if err != nil {
		println("close error")
		panic(err)
	}

	println("Server exited successfully")
}
