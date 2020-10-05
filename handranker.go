package main

import "sort"

type HandRanker struct {
	Board *Board
}

func (h *HandRanker) Straight(playerHand Hand) (int, bool) {
	cards := mergeCards(h.Board.Cards, playerHand.Cards)

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

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

func mergeCards(a, b []Card) []Card {
	cards := a
	cards = append(a, b...)
	return cards
}
