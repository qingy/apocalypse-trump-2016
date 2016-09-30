package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// post the request to the target url
// TODO: check/handle HTTP response code
func postRequest(url string, postBody string) ([]byte, error) {
	var reqBytes = []byte(postBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error posting to %s: %s", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body to %s: %s", postBody, err)
	}

	return respBytes, nil
}

// post JSON to a url
// TODO: check/handle HTTP response code
func postJSON(url string, request interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling request: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error posting to %s: %s", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body to %s: %s", url, err)
	}

	return respBytes, nil
}
