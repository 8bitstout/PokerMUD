package pokermud

import "fmt"

const (
	HIGH_CARD = iota
	ONE_PAIR
	TWO_PAIR
	THREE_OF_A_KIND
	STRAIGHT
	FLUSH
	FULL_HOUSE
	FOUR_OF_A_KIND
	STRAIGHT_FLUSH
)

type Hand struct {
	Cards []Card
	Name  string
	Rank  int
	Value int
}

func (h *Hand) GetEachCard() (Card, Card) {
	return h.Cards[0], h.Cards[1]
}

func (h *Hand) GetHighCard() int {
	return MaxInt(h.Cards[0].Value, h.Cards[1].Value)
}

func (h *Hand) GetKicker() int {
	return MinInt(h.Cards[0].Value, h.Cards[1].Value)
}

func (h *Hand) IsQuads() bool {
	return true
}

func (h *Hand) IsPair(b *Board) bool {
	highCard := h.GetHighCard()
	kicker := h.GetKicker()

	if b.ContainsCardValue(highCard) || b.ContainsCardValue(kicker) {
		return true
	}

	return false
}

func (h *Hand) GetRank() int {
	if h.IsQuads() {
		return FOUR_OF_A_KIND
	}
	return ONE_PAIR
}

// GetValue calculates the value of each card in a players hands and
// returns it. If both cards have the same value, we know we have a pair.
// Therefore, we will multiply the value by 2 to denote a made hand.
func (h *Hand) GetValue() int {
	c1, c2 := h.Cards[0].Value, h.Cards[1].Value
	value := c1 + c2
	h.Rank = 1
	if c1 == c2 {
		h.Rank = 2
	}

	return value
}

func (h *Hand) String() string {
	return fmt.Sprint(h.Cards[0].Name, h.Cards[1].Name)
}

func MakeHand() *Hand {
	return &Hand{
		Cards: []Card{},
	}
}
