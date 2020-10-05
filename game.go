package main

import (
	"fmt"
	"time"
)

const (
	POSITION_SMALL_BLIND = 0
	POSITION_BIG_BLIND   = 1
)

type Game struct {
	Table      *Table
	Players    []*Player
	Deck       *Deck
	Board      *Board
	HandRanker *HandRanker
	Round      int
}

func (g *Game) Play() {
	for len(g.Players) >= 2 {
		g.Deck = MakeDeck()
		g.Board = MakeBoard()
		g.HandRanker.Board = g.Board

		fmt.Println("---Dealing Cards---")
		g.DealPlayers()
		// await betting actions
		fmt.Println("---Dealing Flop---")
		g.DealFlop()
		g.Board.DisplayBoard()
		g.AwaitPlayerDecision()
		g.RankPlayerHands()
		fmt.Println("---Dealing Turn---")
		g.DealTurn()
		g.Board.DisplayBoard()
		g.RankPlayerHands()
		fmt.Println("---Dealing River---")
		g.DealRiver()
		g.Board.DisplayBoard()
		g.RankPlayerHands()
		break
	}
}

func (g *Game) AwaitPlayerDecision() {
	startTime := time.Now()
	stopTime := startTime.Add(time.Second * 10)
	fmt.Println(startTime.Unix())
	fmt.Println(stopTime.Unix())

	for time.Now().Unix() != stopTime.Unix() {
		continue
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
		Round:   0,
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
