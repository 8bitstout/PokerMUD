package pokermud

type Position int8

const (
	POSITION_SMALL_BLIND Position = 0
	POSITION_BIG_BLIND   Position = 1
)

func (p Position) String() string {
	return [...]string{"Small blind", "Big blind"}[p]
}
