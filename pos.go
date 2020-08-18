package czzle

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (pos *Pos) GetX() int {
	if pos == nil {
		return 0
	}
	return pos.X
}

func (pos *Pos) SetX(x int) {
	if pos == nil {
		return
	}
	pos.X = x
}

func (pos *Pos) GetY() int {
	if pos == nil {
		return 0
	}
	return pos.Y
}

func (pos *Pos) SetY(y int) {
	if pos == nil {
		return
	}
	pos.Y = y
}
