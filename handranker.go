package main

import "sort"

type HandRanker struct {
	Board *Board
}

func (h *HandRanker) Straight(playerHand Hand) (int, bool) {
	cards := mergeCards(h.Board.Cards, playerHand.Cards)
	sortCardsByValue(cards)

	count := 0
	highestValue := 0

	for i := 1; i < len(cards); i++ {
		previous := cards[i-1].Value
		if previous+1 != cards[i].Value {
			count++
			highestValue = cards[i].Value
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

	if len(suiteValues) < 5 {
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

func mergeCards(a, b []Card) []Card {
	cards := a
	cards = append(a, b...)
	return cards
}

func sortCardsByValue(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})
}
