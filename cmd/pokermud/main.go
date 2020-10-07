package main

import "github.com/8bitstout/pokermud"

func main() {
	var players []*pokermud.Player

	names := []string{"WCGRider", "OMGClayAiken", "Sauce123", "Ben86"}

	for _, name := range names {
		p := pokermud.MakePlayer(name)
		players = append(players, p)
	}
	g := pokermud.MakeGame(players)
	g.Start()

}
