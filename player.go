package pokermud

import "net"

const (
	ACTION_BET = iota
	ACTION_CHECK
	ACTION_FOLD
)

type Player struct {
	Name            string
	ID              int
	Hand            *Hand
	Value           int
	Chips           int
	Connection      net.Conn
	IsActive        bool
	IsAuthenticated bool
}

func (p *Player) AddCard(c Card) {
	p.Hand.Cards = append(p.Hand.Cards, c)
}

func (p *Player) GiveBigBlind() int {
	return 2
}

func MakePlayer(name string, conn net.Conn) *Player {
	return &Player{
		Name:            name,
		Value:           200,
		Hand:            MakeHand(),
		Connection:      conn,
		IsActive:        true,
		IsAuthenticated: false,
	}
}
