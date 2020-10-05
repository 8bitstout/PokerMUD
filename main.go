package main

const (
	PREFLOP = iota
	FLOP
	TURN
	RIVER
)

type Table struct {
	MaxSeats int
	Players  []Player
	Name     string
}

func main() {
	var players []*Player

	names := []string{"WCGRider", "OMGClayAiken", "Sauce123", "Ben86"}

	for _, name := range names {
		p := MakePlayer(name)
		players = append(players, p)
	}
	g := MakeGame(players)
	g.Play()

}

func GetRankName(rank int) string {
	ranks := map[int]string{
		HIGH_CARD:       "High Card",
		ONE_PAIR:        "Pair",
		TWO_PAIR:        "Two Pair",
		THREE_OF_A_KIND: "Three of a Kind",
		STRAIGHT:        "Straight",
		FLUSH:           "Flush",
		FULL_HOUSE:      "Full House",
		FOUR_OF_A_KIND:  "Four of a Kind",
		STRAIGHT_FLUSH:  "Straight Flush",
	}

	return ranks[rank]
}
func MaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
