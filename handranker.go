package main

import (
	"sort"
)

type HandRanker struct {
	Board *Board
}

func (h *HandRanker) Rank(playerHand Hand) (int, int) {
	if maxVal, ok := h.StraightFlush(playerHand); ok {
		return STRAIGHT_FLUSH, maxVal
	}
	if maxVal, ok := h.FourOfAKind(playerHand); ok {
		return FOUR_OF_A_KIND, maxVal
	}
	if maxVal, ok := h.FullHouse(playerHand); ok {
		return FULL_HOUSE, maxVal
	}
	if maxVal, ok := h.Flush(playerHand); ok {
		return FLUSH, maxVal
	}
	if maxVal, ok := h.Straight(playerHand); ok {
		return STRAIGHT, maxVal
	}
	if maxVal, ok := h.ThreeOfAKind(playerHand); ok {
		return THREE_OF_A_KIND, maxVal
	}
	if maxVal, ok := h.TwoPair(playerHand); ok {
		return TWO_PAIR, maxVal
	}
	if maxVal, ok := h.Pair(playerHand); ok {
		return ONE_PAIR, maxVal
	}
	return HIGH_CARD, h.GetHighCardValue(playerHand)
}

func (h *HandRanker) GetHighCardValue(playerHand Hand) int {
	cards := mergeCards(h.Board.Cards, playerHand.Cards)
	sortCardsByValue(cards)
	return cards[len(cards)-1].Value
}

func (h *HandRanker) FullHouse(playerHand Hand) (int, bool) {
	var foundPair, foundSet bool
	var maxValue int
	cardCount := make(map[int]int)
	cards := mergeCards(h.Board.Cards, playerHand.Cards)

	for _, card := range cards {
		cardCount[card.Value] += 1
	}

	for value, count := range cardCount {
		if count == 2 {
			foundPair = true
		}

		if count == 3 {
			foundSet = true
			maxValue = value
		}
	}

	if foundPair && foundSet {
		return maxValue, true
	}

	return maxValue, false
}

func (h *HandRanker) StraightFlush(playerHand Hand) (int, bool) {
	cards := mergeCards(h.Board.Cards, playerHand.Cards)
	sortCardsByValue(cards)

	count := 0
	highestValue := 0

	for i := 1; i < len(cards); i++ {
		previous := cards[i-1]
		if previous.Value+1 == cards[i].Value && previous.Suite == cards[i].Suite {
			count++
			highestValue = cards[i].Value
		} else {
			count = 0
		}
	}

	return highestValue, count >= 4
}

func (h *HandRanker) Straight(playerHand Hand) (int, bool) {
	cards := mergeCards(h.Board.Cards, playerHand.Cards)
	sortCardsByValue(cards)

	count := 0
	highestValue := 0

	for i := 1; i < len(cards); i++ {
		previous := cards[i-1].Value
		if previous+1 == cards[i].Value {
			count++
			highestValue = cards[i].Value
		} else {
			count = 0
		}
	}
	return highestValue, count >= 4
}

func (h *HandRanker) Flush(playerHand Hand) (int, bool) {
	var cardToUse Card
	c1, c2 := playerHand.GetEachCard()
	handIsDoubleSuited := c1.Suite == c2.Suite

	if handIsDoubleSuited || h.Board.GetSuiteCount(c1.Suite) >= 3 {
		cardToUse = c1
	}

	if h.Board.GetSuiteCount(c2.Suite) >= 3 {
		cardToUse = c2
	}

	suiteValues := h.Board.Suites[cardToUse.Suite]
	suiteValues = append(suiteValues, cardToUse.Value)

	if handIsDoubleSuited && len(suiteValues) < 4 || !handIsDoubleSuited && len(suiteValues) < 5 {
		return -1, false
	}

	maxValue := -1

	for _, value := range suiteValues {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue, true
}

func (h *HandRanker) FourOfAKind(playerHand Hand) (int, bool) {
	c1, c2 := playerHand.GetEachCard()
	handIsPair := c1.Value == c2.Value

	if handIsPair && len(h.Board.Values[c1.Value]) == 2 {
		return c1.Value, true
	}

	if len(h.Board.Values[c1.Value]) == 3 {
		return c1.Value, true
	}

	if len(h.Board.Values[c2.Value]) == 3 {
		return c2.Value, true
	}

	return -1, false
}

func (h *HandRanker) ThreeOfAKind(playerHand Hand) (int, bool) {
	c1, c2 := playerHand.GetEachCard()
	handIsPair := c1.Value == c2.Value

	if handIsPair && len(h.Board.Values[c1.Value]) == 1 {
		return c1.Value, true
	}

	if len(h.Board.Values[c1.Value]) == 2 {
		return c1.Value, true
	}

	if len(h.Board.Values[c2.Value]) == 2 {
		return c2.Value, true
	}

	return -1, false
}

func (h *HandRanker) TwoPair(playerHand Hand) (int, bool) {
	cardCount := make(map[int]int)
	cards := mergeCards(h.Board.Cards, playerHand.Cards)

	for _, card := range cards {
		cardCount[card.Value] += 1
	}

	maxValue := 0
	pairs := 0

	for value, count := range cardCount {
		if count == 2 {
			pairs++
			if value > maxValue {
				maxValue = value
			}
		}
	}

	return maxValue, pairs >= 2
}

func (h *HandRanker) Pair(playerHand Hand) (int, bool) {
	cardCount := make(map[int]int)
	cards := mergeCards(h.Board.Cards, playerHand.Cards)

	for _, card := range cards {
		cardCount[card.Value] += 1
	}

	maxValue := 0
	pairs := 0

	for value, count := range cardCount {
		if count == 2 {
			pairs++
			if value > maxValue {
				maxValue = value
			}
		}
	}

	return maxValue, pairs >= 1
}

func mergeCards(a, b []Card) []Card {
	var cards []Card
	cards = append(cards, a...)
	cards = append(cards, b...)
	return cards
}

func sortCardsByValue(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})
}
