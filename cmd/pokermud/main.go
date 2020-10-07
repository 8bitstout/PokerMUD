package main

import (
	"fmt"
	"github.com/8bitstout/pokermud"
	client2 "github.com/8bitstout/pokermud/client"
	server2 "github.com/8bitstout/pokermud/server"
	"os"
)

func main() {
	var players []*pokermud.Player
	arguments := os.Args

	names := []string{"WCGRider", "OMGClayAiken", "Sauce123", "Ben86"}

	for _, name := range names {
		p := pokermud.MakePlayer(name)
		players = append(players, p)
	}
	//g := pokermud.MakeGame(players)
	//g.Start()

	command := arguments[1]

	fmt.Println(arguments[1])

	if command == "server" {
		s := server2.MakeServer("127.0.0.1:1234")
		s.Start()
	}

	if command == "client" {
		client := client2.MakeClient("127.0.0.1:1234")
		client.Connect()
	}

}
