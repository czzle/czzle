package czzle

type Solution struct {
	Token   string     `json:"token"`
	Actions ActionList `json:"actions"`
}

func (s *Solution) GetToken() string {
	if s == nil {
		return ""
	}
	return s.Token
}

func (s *Solution) SetToken(tkn string) {
	if s == nil {
		return
	}
	s.Token = tkn
}

func (s *Solution) GetActions() ActionList {
	if s == nil {
		return nil
	}
	return s.Actions
}

func (s *Solution) SetActions(actions ActionList) {
	if s == nil {
		return
	}
	s.Actions = actions
}
