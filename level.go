package czzle

type Level string

const (
	UnknownLevel Level = "unknown"
	None         Level = "none"
	Easy         Level = "easy"
	Medium       Level = "medium"
	Hard         Level = "hard"
)

func (lvl Level) UInt8() uint8 {
	switch lvl {
	case None:
		return 1
	case Easy:
		return 2
	case Medium:
		return 3
	case Hard:
		return 4
	default:
		return 0
	}
}
