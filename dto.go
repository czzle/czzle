package czzle

type BeginReq struct {
	Client *ClientInfo `json:"client"`
}

func (req *BeginReq) GetClient() *ClientInfo {
	if req == nil {
		return nil
	}
	return req.Client
}

func (req *BeginReq) SetClient(client *ClientInfo) {
	if req == nil {
		return
	}
	req.Client = client
}

type BeginRes struct {
	Puzzle *Puzzle `json:"puzzle"`
}

func (res *BeginRes) GetPuzzle() *Puzzle {
	if res == nil {
		return nil
	}
	return res.Puzzle
}

func (res *BeginRes) SetPuzzle(puzzle *Puzzle) {
	if res == nil {
		return
	}
	res.Puzzle = puzzle
}

type SolveReq struct {
	Solution *Solution `json:"solution"`
}

func (req *SolveReq) GetSolution() *Solution {
	if req == nil {
		return nil
	}
	return req.Solution
}

func (req *SolveReq) SetSolution(solution *Solution) {
	if req == nil {
		return
	}
	req.Solution = solution
}

type SolveRes struct {
	Results *Results `json:"results"`
}

func (res *SolveRes) GetResults() *Results {
	if res == nil {
		return nil
	}
	return res.Results
}

func (res *SolveRes) SetResults(results *Results) {
	if res == nil {
		return
	}
	res.Results = results
}

type ValidateReq struct {
	AccessToken string `json:"access_token"`
}

func (req *ValidateReq) GetAccessToken() string {
	if req == nil {
		return ""
	}
	return req.AccessToken
}

func (req *ValidateReq) SetAccessToken(tkn string) {
	if req == nil {
		return
	}
	req.AccessToken = tkn
}

type ValidateRes struct {
	OK bool `json:"ok"`
}

func (res *ValidateRes) IsOK() bool {
	if res == nil {
		return false
	}
	return res.OK
}

func (res *ValidateRes) SetOK(ok bool) {
	if res == nil {
		return
	}
	res.OK = ok
}
