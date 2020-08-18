package czzle

import (
	"encoding/json"
)

type ActionData interface {
	actionData()
}

func (*FlipActionData) actionData() {}

type ActionType string

const (
	UnknownAction ActionType = "unknown"
	BeginAction   ActionType = "begin"
	FlipAction    ActionType = "flip"
	ConfirmAction ActionType = "confirm"
)

type ActionTypeMap map[ActionType]bool

func (m ActionTypeMap) Has(tt ActionType) bool {
	ok, has := m[tt]
	return has && ok
}

var AllowedActions = ActionTypeMap{
	BeginAction:   true,
	FlipAction:    true,
	ConfirmAction: true,
}

type Action struct {
	Type ActionType
	Time int64
	Data ActionData
}

type ActionList []*Action

func (list ActionList) Len() int {
	return len(list)
}

func (list ActionList) Less(i, j int) bool {
	return list[i].GetTime() < list[j].GetTime()
}

func (list ActionList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type jsonAction struct {
	Type ActionType      `json:"type"`
	Time int64           `json:"time"`
	Data json.RawMessage `json:"data,omitempty"`
}

func (a Action) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(a.Data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonAction{
		Type: a.Type,
		Time: a.Time,
		Data: data,
	})
}

func (a *Action) UnmarshalJSON(src []byte) error {
	var raw jsonAction
	err := json.Unmarshal(src, &raw)
	if err != nil {
		return err
	}
	var data ActionData
	switch raw.Type {
	case FlipAction:
		data = new(FlipActionData)
	}
	if data != nil {
		err = json.Unmarshal(raw.Data, data)
	}
	*a = Action{
		Type: raw.Type,
		Time: raw.Time,
		Data: data,
	}
	return nil
}

func (a *Action) GetType() ActionType {
	if a == nil {
		return UnknownAction
	}
	return a.Type
}

func (a *Action) SetType(typ ActionType) {
	if a == nil {
		return
	}
	a.Type = typ
}

func (a *Action) GetTime() int64 {
	if a == nil {
		return 0
	}
	return a.Time
}

func (a *Action) SetTime(ts int64) {
	if a == nil {
		return
	}
	a.Time = ts
}

func (a *Action) GetData() ActionData {
	if a == nil {
		return nil
	}
	return a.Data
}

func (a *Action) SetData(data ActionData) {
	if a == nil {
		return
	}
	a.Data = data
}

type FlipActionData struct {
	Pos *Pos `json:"pos"`
}

func (data *FlipActionData) GetPos() *Pos {
	if data == nil {
		return nil
	}
	return data.Pos
}

func (data *FlipActionData) SetPos(pos *Pos) {
	if data == nil {
		return
	}
	data.Pos = pos
}
