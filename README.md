# jsonrpc2

```golang
package main

import (
	"errors"
	"./jsonrpc2"
	"net/http"
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

func myFn(r *http.Request) bool {
	return true
}

func main(){
	server := jsonrpc2.NewServer("/jsonrpc", "0.0.0.0", "8008", myFn)
	server.RegisterFunc("echo", echo)

	server.Start()
}
```
