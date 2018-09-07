package jsonrpc2

import (
	"encoding/json"
)

type Error struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func NewError(id interface{}, code int16, message string, data string) []byte {
	resp := Response{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  nil,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	res, _ := json.Marshal(&resp)
	return res
}
