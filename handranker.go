package main

import "sort"

type HandRanker struct {
	Board *Board
}

func (h *HandRanker) IsStraight(playerHand Hand) (bool, int) {
	cards := h.Board.Cards
	cards = append(cards, playerHand.Cards...)

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

	return count >= 4, highestValue
}
