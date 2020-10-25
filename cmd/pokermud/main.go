package main

import (
	"fmt"
	"github.com/8bitstout/pokermud/tcp"
	"os"
)

func main() {
	arguments := os.Args
	command := arguments[1]
	fmt.Println(arguments[1])

	os.Setenv("ENABLE_LOGGING", "0")

	if len(arguments) > 2 && arguments[2] == "--debug" {
		fmt.Println("Logging enabled")
		os.Setenv("ENABLE_LOGGING", "1")
	}

	if command == "server" {
		s := tcp.MakeServer("127.0.0.1:1234")
		s.Start()
	}

	if command == "client" {
		client := tcp.MakeClient("127.0.0.1:1234")
		client.Connect()
	}

}
