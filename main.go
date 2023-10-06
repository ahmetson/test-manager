package main

import (
	"fmt"
	clientConfig "github.com/ahmetson/client-lib/config"
	"github.com/ahmetson/datatype-lib/data_type/key_value"
	"github.com/ahmetson/datatype-lib/message"
	"github.com/ahmetson/os-lib/arg"
	"github.com/pebbe/zmq4"
	"os"
	"time"
)

const PortFlag = "port"
const ManagerPortFlag = "manager-port"
const IdFlag = "id"
const UrlFlag = "url"
const ParentFlag = "parent"

//const ParentFlag = "parent"

func runHandler(port string) {
	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		panic(err)
	}

	err = socket.Bind(fmt.Sprintf("tcp://*:%s", port))
	if err != nil {
		println("runHandler: bind error")
		panic(err)
	}
	defer func() {
		err := socket.Close()
		if err != nil {
			panic(err)
		}
	}()

	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)

	currentTime := time.Now()

	fmt.Printf("runHandler: current time: %v, server running on port %s\n", currentTime, port)

	for {
		polled, err := poller.Poll(time.Millisecond)
		if err != nil {
			panic(err)
		}

		if len(polled) == 0 {
			continue
		}

		raw, err := socket.RecvMessage(0)
		if err != nil {
			panic(err)
		}

		req, err := message.NewReq(raw)
		if err != nil {
			panic(err)
		}

		reply := req.Ok(key_value.New().Set("command", req.CommandName()))
		replyStr := reply.String()
		_, err = socket.SendMessage(replyStr)
		if err != nil {
			panic(err)
		}
	}
}

func runManager(port string) {
	closeFlag := false

	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		panic(err)
	}

	err = socket.Bind(fmt.Sprintf("tcp://*:%s", port))
	if err != nil {
		println("runManager: bind error")
		panic(err)
	}

	poller := zmq4.NewPoller()
	poller.Add(socket, zmq4.POLLIN)

	currentTime := time.Now()

	fmt.Printf("runManager: current time: %v, server running on port %s\n", currentTime, port)

	for {
		if closeFlag {
			break
		}

		polled, err := poller.Poll(time.Millisecond)
		if err != nil {
			panic(err)
		}

		if len(polled) == 0 {
			continue
		}

		raw, err := socket.RecvMessage(0)
		if err != nil {
			panic(err)
		}

		req, err := message.NewReq(raw)
		if err != nil {
			panic(err)
		}

		if req.CommandName() == "close" {
			closeFlag = true
		}

		reply := req.Ok(key_value.New().Set("command", req.CommandName()))
		replyStr := reply.String()
		_, err = socket.SendMessage(replyStr)
		if err != nil {
			panic(err)
		}
	}

	err = socket.Close()
	if err != nil {
		println("runManager: close error")
		panic(err)
	}

	fmt.Printf("Proxy exited successfully: %v, now: %v\n", currentTime, time.Now())
	os.Exit(0)
}

// For N seconds, it imitates the SDS handler as busy.
// For testing dep_manager.Running() command.
func main() {
	port := "6000" // default port
	if arg.FlagExist(PortFlag) {
		port = arg.FlagValue(PortFlag)
	}
	managerPort := "6001"
	if arg.FlagExist(ManagerPortFlag) {
		managerPort = arg.FlagValue(ManagerPortFlag)
	}
	if !arg.FlagExist(IdFlag) {
		panic("no id flag")
	}
	if !arg.FlagExist(UrlFlag) {
		panic("no url flag")
	}
	if !arg.FlagExist(ParentFlag) {
		panic("no parent flag")
	}
	id := arg.FlagValue(IdFlag)
	url := arg.FlagValue(UrlFlag)
	parentRaw := arg.FlagValue(ParentFlag)
	parentKv, err := key_value.NewFromString(parentRaw)
	if err != nil {
		panic(err)
	}

	var parentConfig clientConfig.Client
	err = parentKv.Interface(&parentConfig)
	if err != nil {
		panic(err)
	}

	//parentConfig.UrlFunc(clientConfig.Url)
	//parent, err := client.New(&parentConfig)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Printf("proxy imitating with id='%s', url='%s', port='%s', managerPort='%s'\n", id, url, port, managerPort)

	fmt.Printf("proxy generated a configuration...notify the parent\n")
	fmt.Printf("fetch from the parent the proxy chain\n")
	fmt.Printf("fetch from the parent the proxy units\n")

	time.Sleep(time.Second)

	fmt.Printf("proxy imitated that all of it's proxies are set\n")
	fmt.Printf("notify parent that PROXY READY")

	go runHandler(port)
	runManager(managerPort)
}
