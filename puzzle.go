package czzle

import (
	"time"
)

type Puzzle struct {
	Level     Level       `json:"level"`
	Token     string      `json:"token"`
	Client    *ClientInfo `json:"client"`
	ExpiresAt int64       `json:"expires_at"`
	IssuedAt  int64       `json:"issued_at"`
	TileMap   *TileMap    `json:"tile_map"`
}

func (p *Puzzle) GetLevel() Level {
	if p == nil {
		return UnknownLevel
	}
	return p.Level
}

func (p *Puzzle) SetLevel(lvl Level) {
	if p == nil {
		return
	}
	p.Level = lvl
}

func (p *Puzzle) GetToken() string {
	if p == nil {
		return ""
	}
	return p.Token
}

func (p *Puzzle) SetToken(token string) {
	if p == nil {
		return
	}
	p.Token = token
}

func (p *Puzzle) GetClient() *ClientInfo {
	if p == nil {
		return nil
	}
	return p.Client
}

func (p *Puzzle) SetClient(client *ClientInfo) {
	if p == nil {
		return
	}
	p.Client = client
}

func (p *Puzzle) GetExpirationTime() int64 {
	if p == nil {
		return 0
	}
	return p.ExpiresAt
}

func (p *Puzzle) IsExpired() bool {
	if p == nil {
		return false
	}
	sec := p.ExpiresAt / 1000
	nsec := (p.ExpiresAt % 1000) * 1000000
	exp := time.Unix(sec, nsec)
	return time.Now().After(exp)
}

func (p *Puzzle) SetExpirationTime(ts int64) {
	if p == nil {
		return
	}
	p.ExpiresAt = ts
}

func (p *Puzzle) GetIssuedTime() int64 {
	if p == nil {
		return 0
	}
	return p.IssuedAt
}

func (p *Puzzle) SetIssuedTime(ts int64) {
	if p == nil {
		return
	}
	p.IssuedAt = ts
}

func (p *Puzzle) GetTileMap() *TileMap {
	if p == nil {
		return nil
	}
	return p.TileMap
}

func (p *Puzzle) SetTileMap(tm *TileMap) {
	if p == nil {
		return
	}
	p.TileMap = tm
}
