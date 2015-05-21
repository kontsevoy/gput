package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	identityUrl = "https://identity.api.rackspacecloud.com/v2.0/tokens"
	apiKey      = "f0a0d5090040ff0650a76fdf751537c7"
	authRequest = `{ "auth": {
				"RAX-KSKEY:apiKeyCredentials": {
					"apiKey": "%s",
					"username": "ekontsevoy"
				}}}`
)

func authenticate() (authToken string, err error) {
	requestBody := fmt.Sprintf(authRequest, apiKey)
	fmt.Println(requestBody)

	// HTTP POST to auth URL:
	response, err := http.Post(identityUrl, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = errors.New(response.Status)
		return
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	panicIf(err)

	// Parse returned JSON:
	//fmt.Printf("Response.body: %v\nResponse.Status: %v\n",
	//		string(responseBody), response.Status)

	var m map[string]interface{}
	json.Unmarshal(responseBody, &m)

	for key, _ := range m {
		fmt.Println(key)
	}

	access, found := m["access"]
	if !found {
		err = errors.New("Invalid response body")
	}

	m, ok := access.(map[string]interface{})
	if ok {
		printKeys(m)
	}

	return "", err
}

func main() {
	authToken, err := authenticate()
	if err != nil {
		fmt.Printf("%v when trying to authenticate\n", err)
		return
	}
	fmt.Printf("%v\n", authToken)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func printKeys(m map[string]interface{}) {
	for key, _ := range m {
		fmt.Println(key)
	}
}
