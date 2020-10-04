package main

type Player struct {
	Name  string
	Hand  *Hand
	Value int
	Chips int
}

func (p *Player) AddCard(c Card) {
	p.Hand.Cards = append(p.Hand.Cards, c)
}

func (p *Player) GiveBigBlind() int {
	return 2
}

func MakePlayer(name string) *Player {
	return &Player{
		Name:  name,
		Value: 200,
		Hand:  MakeHand(),
	}
}
