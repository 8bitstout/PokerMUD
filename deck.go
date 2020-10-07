package pokermud

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	HEART = iota
	DIAMOND
	SPADE
	CLUB
)

type Card struct {
	Value int
	Suite int
	Name  string
}

func (c *Card) SuiteToString() string {
	switch s := c.Suite; s {
	case HEART:
		return "h"
	case DIAMOND:
		return "d"
	case SPADE:
		return "s"
	case CLUB:
		return "c"
	}

	return ""
}

type Deck struct {
	Cards []Card
}

func (d *Deck) RemoveTopCard() Card {
	c := d.Cards[len(d.Cards)-1]
	d.Cards = d.Cards[:len(d.Cards)-1]

	return c
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
}

func MakeDeck() *Deck {
	suites := [4]int{HEART, DIAMOND, CLUB, SPADE}
	faceCards := map[int]string{
		11: "J",
		12: "Q",
		13: "K",
		14: "A",
	}
	d := &Deck{
		[]Card{},
	}

	for i := 0; i < len(suites); i++ {
		for j := 2; j <= 14; j++ {
			c := &Card{
				Value: j,
				Suite: suites[i],
			}
			c.Name = fmt.Sprint(j, c.SuiteToString())

			if _, ok := faceCards[j]; ok {
				c.Name = fmt.Sprint(faceCards[j], c.SuiteToString())
			}

			d.Cards = append(d.Cards, *c)
		}
	}

	return d
}
