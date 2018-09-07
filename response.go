package jsonrpc2

// {"jsonrpc": "2.0", "result": 19, "id": 1}
type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
