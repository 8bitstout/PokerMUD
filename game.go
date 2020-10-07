package pokermud

import (
	"fmt"
	"log"
	"time"
)

type Game struct {
	Players    []*Player
	Deck       *Deck
	Board      *Board
	HandRanker *HandRanker
	Round      Round
}

func (g *Game) Start() {
	g.SetRound(ROUND_PREFLOP)
	for g.Round != ROUND_WAITING_FOR_PLAYERS {
		g.Deck = MakeDeck()
		g.Board = MakeBoard()
		g.HandRanker.Board = g.Board

		fmt.Println("---Dealing Cards---")
		g.DealPlayers()
		g.ForEachPlayer(g.AwaitPlayerDecision)
		g.RankPlayerHands()

		fmt.Println("---Dealing Flop---")
		g.DealFlop()
		g.Board.DisplayBoard()
		g.ForEachPlayer(g.AwaitPlayerDecision)
		g.RankPlayerHands()

		fmt.Println("---Dealing Turn---")
		g.DealTurn()
		g.Board.DisplayBoard()
		g.RankPlayerHands()
		g.ForEachPlayer(g.AwaitPlayerDecision)

		fmt.Println("---Dealing River---")
		g.DealRiver()
		g.Board.DisplayBoard()
		g.RankPlayerHands()
		g.ForEachPlayer(g.AwaitPlayerDecision)
		break
	}
}

func (g *Game) SetRound(r Round) {
	g.Round = r
}

func (g *Game) ForEachPlayer(f func(p *Player)) {
	for _, player := range g.Players {
		f(player)
	}
}

func (g *Game) AwaitPlayerDecision(player *Player) {
	if !player.IsActive {
		log.Printf(player.Name, "is inactive")
		return
	}
	log.Println(player.Name, "is making a decision")
	playerMadeDecision := false
	startTime := time.Now()
	stopTime := startTime.Add(time.Second * 10)
	fmt.Println(startTime.Unix())
	fmt.Println(stopTime.Unix())

	for time.Now().Unix() != stopTime.Unix() {
		continue
	}

	if !playerMadeDecision {
		log.Printf(player.Name, "failed to make a decision in time")
		player.IsActive = false
	}
}

func (g *Game) GetPlayerInBigBlind() *Player {
	return g.Players[POSITION_BIG_BLIND]
}

func (g *Game) GetPlayerInSmallBlind() *Player {
	return g.Players[POSITION_SMALL_BLIND]
}

func (g *Game) RankPlayerHands() {
	for _, player := range g.Players {
		rank, _ := g.HandRanker.Rank(*player.Hand)
		fmt.Println(player.Name, "has a", GetRankName(rank), player.Hand.ToString())
	}
}

func MakeGame(players []*Player) *Game {
	d := MakeDeck()
	b := MakeBoard()

	g := &Game{
		Deck:    d,
		Board:   b,
		Players: players,
		Round:   ROUND_WAITING_FOR_PLAYERS,
	}
	h := &HandRanker{
		Board: g.Board,
	}
	g.HandRanker = h

	return g
}

func (g *Game) DealPlayers() {
	g.Deck.Shuffle()
	for i := 0; i < 2; i++ {
		for _, player := range g.Players {
			player.AddCard(g.Deck.RemoveTopCard())
		}
	}
}

func (g *Game) DealFlop() {
	g.Deck.RemoveTopCard()
	for i := 0; i < 3; i++ {
		g.Board.AddCard(g.Deck.RemoveTopCard())
	}
}

func (g *Game) DealTurn() {
	g.Deck.RemoveTopCard()
	g.Board.AddCard(g.Deck.RemoveTopCard())
}

func (g *Game) DealRiver() {
	g.Deck.RemoveTopCard()
	g.Board.AddCard(g.Deck.RemoveTopCard())
}

func (g *Game) RotatePlayers() {
	current := g.Players[0]
	previous := current
	for i := 1; i < len(g.Players); i++ {
		current = g.Players[i]
		g.Players[i] = previous
		previous = current
	}
	g.Players[0] = previous
}

func (g *Game) ListPlayers() {
	for _, p := range g.Players {
		fmt.Println(p.Name)
	}
}

func GetRankName(rank int) string {
	ranks := map[int]string{
		HIGH_CARD:       "High Card",
		ONE_PAIR:        "Pair",
		TWO_PAIR:        "Two Pair",
		THREE_OF_A_KIND: "Three of a Kind",
		STRAIGHT:        "Straight",
		FLUSH:           "Flush",
		FULL_HOUSE:      "Full House",
		FOUR_OF_A_KIND:  "Four of a Kind",
		STRAIGHT_FLUSH:  "Straight Flush",
	}

	return ranks[rank]
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
