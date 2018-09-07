package jsonrpc2

import (
	"encoding/json"
)

// {"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}
type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}
