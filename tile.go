package czzle

type TileType string

const (
	UnknownTile TileType = "unknown"
	FrontTile   TileType = "front"
	BackTile    TileType = "back"
)

type Tile struct {
	Type TileType `json:"type"`
	Pos  *Pos     `json:"pos"`
	Data []byte   `json:"data"`
}

func (t *Tile) GetType() TileType {
	if t == nil {
		return UnknownTile
	}
	return t.Type
}

func (t *Tile) SetType(typ TileType) {
	if t == nil {
		return
	}
	t.Type = typ
}

func (t *Tile) GetPos() *Pos {
	if t == nil {
		return nil
	}
	return t.Pos
}

func (t *Tile) SetPos(pos *Pos) {
	if t == nil {
		return
	}
	t.Pos = pos
}

func (t *Tile) GetData() []byte {
	if t == nil {
		return nil
	}
	return t.Data
}

func (t *Tile) SetData(data []byte) {
	if t == nil {
		return
	}
	t.Data = data
}
