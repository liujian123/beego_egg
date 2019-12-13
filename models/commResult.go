package models

type CommResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCommResult() *CommResult {

	return &CommResult{}
}
