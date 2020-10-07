package main

type Round int8

const (
	ROUND_WAITING_FOR_PLAYERS Round = iota
	ROUND_PREFLOP             Round = iota
	ROUND_FLOP                Round = iota
	ROUND_TURN                Round = iota
	ROUND_RIVER               Round = iota
)

func (r Round) String() string {
	return [...]string{
		"Waiting for players to join",
		"Preflop",
		"Flop",
		"Turn",
		"River",
	}[r]
}
