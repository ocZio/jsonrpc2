# jsonrpc2

```golang
package jsonrpc2

import (
	"errors"
)

type echoParams struct {
	Message  *string `json:"message"`
}

func echo(params *echoParams) (interface{}, error) {
	if params.Message == nil {
		return nil, errors.New("missing message")
	}

	return params.Message, nil
}

func StartJSONRPCServer(entryPoint string, ip string, port string) {
	server := NewServer(entryPoint, ip, port)

	server.RegisterFunc("echo", echo)

	server.Start()
}
```
