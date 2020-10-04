package main

import "fmt"

type Board struct {
	Cards  []Card
	Value  int
	Rank   int
	Values [15][]int
	Suites [4][]int
}

func (b *Board) AddCard(c Card) {
	b.Cards = append(b.Cards, c)
	b.AddSuite(c)
	b.AddValue(c)
}

func (b *Board) AddSuite(c Card) {
	b.Suites[c.Suite] = append(b.Suites[c.Suite], c.Value)
}

func (b *Board) AddValue(c Card) {
	b.Values[c.Value] = append(b.Values[c.Value], c.Value)
}

func (b *Board) CalculateRank() {
	if b.HasFourOfAKind() {
		b.Rank = FOUR_OF_A_KIND
	} else if b.HasFlush() {
		b.Rank = FLUSH
		return
	} else if b.HasThreeOfAKind() {
		b.Rank = THREE_OF_A_KIND
	} else if b.HasTwoPair() {
		b.Rank = TWO_PAIR
	} else if b.HasPair() {
		b.Rank = ONE_PAIR
	} else {
		b.Rank = HIGH_CARD
	}

}

func (b *Board) HasFlush() bool {
	if b.GetSuiteCount(HEART) == 5 ||
		b.GetSuiteCount(DIAMOND) == 5 ||
		b.GetSuiteCount(CLUB) == 5 ||
		b.GetSuiteCount(SPADE) == 5 {
		return true
	}

	return false
}

func (b *Board) HasNumberOfValues(n int) bool {
	for _, v := range b.Values {
		if len(v) == n {
			return true
		}
	}
	return false
}

func (b *Board) HasFourOfAKind() bool {
	return b.HasNumberOfValues(4)
}

func (b *Board) HasThreeOfAKind() bool {
	return b.HasNumberOfValues(3)
}

func (b *Board) HasPair() bool {
	return b.HasNumberOfValues(2)
}

func (b *Board) HasTwoPair() bool {
	pairCount := 0
	for _, v := range b.Values {
		if len(v) == 2 {
			pairCount++
		}
	}
	return pairCount >= 2
}

func (b *Board) DisplayBoard() {
	fmt.Println("---Board---")
	c1, c2, c3 := b.Cards[0], b.Cards[1], b.Cards[2]
	fmt.Println(c1.Name, c2.Name, c3.Name)
}

func (b *Board) GetSuiteCount(s int) int {
	return len(b.Suites[s])
}

func (b *Board) GetValueCount(v int) int {
	return len(b.Suites[v])
}

func (b *Board) ContainsSuite(s int) bool {
	if len(b.Suites[s]) == 0 {
		return false
	}
	return true
}

func (b *Board) ContainsCardValue(v int) bool {
	if len(b.Values[v]) == 0 {
		return false
	}
	return true
}

func MakeBoard() *Board {
	var cards []Card
	return &Board{
		Cards:  cards,
		Value:  0,
		Values: [15][]int{},
		Suites: [4][]int{},
	}
}
