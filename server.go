package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type handler struct {
	method reflect.Value
	params interface{}
}

type server struct {
	handlers   map[string]handler
	ip_port    string
	entrypoint string
}

func NewServer(entrypoint, ip, port string) *server {
	return &server{
		handlers:   make(map[string]handler),
		ip_port:    fmt.Sprintf("%s:%s", ip, port),
		entrypoint: entrypoint,
	}
}

func jsonrpc(rpcserver *server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: Method needs to be POST, you have used %s", r.Method)
		return
	}

	if _, ok := r.Header["Content-Type"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Content-Type not defined")
		return
	}

	if r.Header["Content-Type"][0] != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Content-Type needs to be application/json, you have used %s", r.Header["Content-Type"][0])
		return
	}

	// read the body content of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	// try to parse the request
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError("0", -32711, "Parse error", ""))
		return
	}

	// as per jsonrpc specification, INT, STRING, NULL
	switch req.Id.(type) {
	case float64:
		req.Id = int64(req.Id.(float64))
	case string:
		req.Id = req.Id.(string)
	default:
		req.Id = nil
	}

	// try to find a method to call
	handler, ok := rpcserver.handlers[req.Method]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(NewError(req.Id, -32601, "Method not found", ""))
		return
	}

	// decode parameters
	params, err := handler.decodeParams(req.Params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(req.Id, -32602, "Invalid method parameter(s)", err.Error()))
		return
	}

	// call the method and see if any errors...
	result, err := handler.call(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(req.Id, -32602, "Invalid method parameter(s)", err.Error()))
		return
	}

	b, _ := json.Marshal(Response{"2.0", req.Id, result, nil})
	w.Write(b)
}

func (s *server) Start() {
	// main entry point for the http, everything else will yield 404.
	http.HandleFunc(s.entrypoint, func(w http.ResponseWriter, r *http.Request) {
		jsonrpc(s, w, r)
	})

	// start the server
	log.Println("Starting JSONRPC server...", s.ip_port)
	err := http.ListenAndServe(s.ip_port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *server) addHandler(method string, fn reflect.Value) {
	if _, found := s.handlers[method]; !found {
		ft := fn.Type()
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.NumIn() != 1 || ft.NumOut() != 2 || !ft.In(0).Implements(reflect.TypeOf((*interface{})(nil)).Elem()) {
			panic(fmt.Sprintf("Method '%s' will not be registered, invalid signature", method))
		}

		if !ft.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			panic(fmt.Sprintf("Method '%s' will not be registered, invalid return signature", method))
		}

		params := ft.In(0)

		if params.Kind() == reflect.Ptr {
			params = params.Elem()
		}

		s.handlers[method] = handler{
			method: fn,
			params: reflect.New(params).Interface().(interface{}),
		}
		log.Printf("Registering '%s'...", method)
	} else {
		log.Printf("Method '%s' already exists, ignoring...", method)
	}
}

func (s *server) RegisterFunc(method string, fn interface{}) {
	s.addHandler(method, reflect.ValueOf(fn))
}

func (h *handler) decodeParams(message json.RawMessage) (interface{}, error) {
	params := reflect.New(reflect.TypeOf(h.params).Elem()).Interface()
	if err := json.Unmarshal(message, &params); err != nil {
		return nil, err
	}
	return params.(interface{}), nil
}

func (h *handler) call(params interface{}) (interface{}, error) {
	result := h.method.Call([]reflect.Value{reflect.ValueOf(params)})
	if result[1].IsNil() {
		if result[0].IsNil() {
			return nil, nil
		}
		return result[0].Interface(), nil
	}
	return nil, result[1].Interface().(error)
}
