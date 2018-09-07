package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type client struct {
	uri  string
	host string
}

func (c *client) MakeRequest(method string, params interface{}) (interface{}, error) {
	p, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	jsonreq := &Request{"2.0", "1", method, p}
	buf, err := json.Marshal(jsonreq)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", c.uri, body)
	req.Header.Add("Content-Type", "application/json")
	req.Host = c.host

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var parsed interface{}
	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return nil, err
	}

	// TODO: remove debugging
	fmt.Println(string(response))

	error_data, have_error := parsed.(map[string]interface{})["error"]
	if have_error {
		return nil, errors.New(error_data.(map[string]interface{})["message"].(string))
	}

	return parsed.(map[string]interface{})["result"], nil
}

func NewClient(uri string) (*client, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return &client{uri: uri, host: u.Hostname()}, nil
}
