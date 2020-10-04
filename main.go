package main

import (
	"fmt"
)

const (
	PREFLOP = iota
	FLOP
	TURN
	RIVER
)

type Table struct {
	MaxSeats int
	Players  []Player
	Name     string
}

type Game struct {
	Table   *Table
	Players []*Player
	Deck    *Deck
	Round   int
}

func main() {
	var players []*Player
	d := MakeDeck()
	b := MakeBoard()

	names := []string{"WCGRider", "OMGClayAiken", "Sauce123", "Ben86"}

	for _, name := range names {
		p := MakePlayer(name)
		players = append(players, p)
	}

	fmt.Println("---Dealing Cards---")

	deal(b, d, players)

	for _, player := range players {
		fmt.Printf("%s: %s - Kicker %d\n", player.Name, player.Hand.ToString(), player.Hand.GetKicker())
	}

	b.DisplayBoard()

	for _, player := range players {
		if player.Hand.IsPair(b) {
			fmt.Println(player.Name, "has a pair")
		}
	}

	fmt.Println("Board contains", b.GetSuiteCount(DIAMOND), "diamonds")
	b.CalculateRank()
	fmt.Println("Board Rank: ", b.Rank)

}

func deal(board *Board, deck *Deck, players []*Player) {
	deck.Shuffle()
	for i := 0; i < 2; i++ {
		for _, player := range players {
			player.AddCard(deck.RemoveTopCard())
		}
	}

	deck.RemoveTopCard()

	for i := 0; i < 3; i++ {
		board.AddCard(deck.RemoveTopCard())
	}
}

func rotatePlayers(players []*Player) {
	current := players[0]
	previous := current
	for i := 1; i < len(players); i++ {
		current = players[i]
		players[i] = previous
		previous = current
	}
	players[0] = previous
}

func listPlayers(players []*Player) {
	for _, p := range players {
		fmt.Println(p.Name)
	}
}

func rankCards(cards []Card) int {
	// iterate through and check for highest card and any duplicates
	//valFrequencies := make(map[int]int)
	//suiteFrequencies := make(map[string]int)

	maxVal := 0
	for _, card := range cards {
		//valFrequencies[card.Value] += 1
		//suiteFrequencies[card.Suite] += 1
		maxVal = MaxInt(maxVal, card.Value)
	}

	return maxVal
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
