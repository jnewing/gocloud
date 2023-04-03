package cloudflare

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func APICall(method string, uri string, email string, key string, params *bytes.Buffer) []byte {
	client := &http.Client{}

	var payload io.Reader

	if params != nil {
		payload = params
	}

	req, _ := http.NewRequest(method, uri, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", key)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error while sending API request.")
		return nil
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return resp_body
}
