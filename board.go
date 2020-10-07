package pokermud

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

func (b *Board) HasNumberOfValues(n int) bool {
	for _, v := range b.Values {
		if len(v) == n {
			return true
		}
	}
	return false
}

func (b *Board) DisplayBoard() {
	fmt.Println("---Board---")
	cards := ""
	for _, card := range b.Cards {
		cards += fmt.Sprint(card.Name, " ")
	}
	fmt.Println(cards)
}

func (b *Board) GetSuiteCount(s Suite) int {
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
