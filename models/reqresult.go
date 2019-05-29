package models

const (
	SUCCEED = iota
	InternalErr
	InvalidArg
)

var (
	CodeMessage = map[int]string{
		SUCCEED:     "succeed",
		InternalErr: "internal error",
		InvalidArg:  "invalid argument",
	}
)

type OperationResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type QueryResult struct {
	OperationResult
	Value interface{} `json:"value"`
}

type LoginResult struct {
	OperationResult
	Token string `json:"token"`
}

func NewOperationResult(code int) OperationResult {
	return OperationResult{
		code,
		CodeMessage[code],
	}
}
