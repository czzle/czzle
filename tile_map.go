package czzle

type TileMap struct {
	Size  int     `json:"size"`
	Tiles []*Tile `json:"tiles"`
}

func (tm *TileMap) GetSize() int {
	if tm == nil {
		return 0
	}
	return tm.Size
}

func (tm *TileMap) SetSize(size int) {
	if tm == nil {
		return
	}
	tm.Size = size
}

func (tm *TileMap) GetTiles() []*Tile {
	if tm == nil {
		return nil
	}
	return tm.Tiles
}

func (tm *TileMap) SetTiles(tiles []*Tile) {
	if tm == nil {
		return
	}
	tm.Tiles = tiles
}
