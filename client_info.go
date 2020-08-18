package czzle

import "github.com/czzle/czzle/pkg/uuid"

type ClientInfo struct {
	ID        uuid.UUID `json:"id"`
	IP        string    `json:"ip"`
	Time      int64     `json:"time"`
	UserAgent string    `json:"user_agent"`
}

func (info *ClientInfo) GetID() uuid.UUID {
	if info == nil {
		return uuid.Null()
	}
	return info.ID
}

func (info *ClientInfo) SetID(id uuid.UUID) {
	if info == nil {
		return
	}
	info.ID = id
}

func (info *ClientInfo) GetIP() string {
	if info == nil {
		return ""
	}
	return info.IP
}

func (info *ClientInfo) SetIP(ip string) {
	if info == nil {
		return
	}
	info.IP = ip
}

func (info *ClientInfo) GetTime() int64 {
	if info == nil {
		return 0
	}
	return info.Time
}

func (info *ClientInfo) SetTime(ts int64) {
	if info == nil {
		return
	}
	info.Time = ts
}

func (info *ClientInfo) GetUserAgent() string {
	if info == nil {
		return ""
	}
	return info.UserAgent
}

func (info *ClientInfo) SetUserAgent(agent string) {
	if info == nil {
		return
	}
	info.UserAgent = agent
}
