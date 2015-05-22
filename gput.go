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

	value, _ := findJsonValue("access/serviceCatalog/endpoints/publicURL", m)
	fmt.Println(strings.Join(value, "\n"))

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

type jsonMap map[string]interface{}

func findJsonValue(searchKey string, m jsonMap) (foundValues []string, found bool) {
	var _search func(jsonMap, []string)
	_search = func(m jsonMap, path []string) {
		keys := mapKeys(m)
		for _, key := range keys {
			value := m[key]
			switch value.(type) {
			case string:
				p := append(path, key)
				fmt.Printf("%v\n", strings.Join(p, "/"))
				if searchKey == strings.Join(p, "/") {
					foundValues = append(foundValues, value.(string))
					found = true
				}
			case map[string]interface{}:
				_search(value.(map[string]interface{}), append(path, key))

			case []interface{}:
				array, _ := value.([]interface{})
				for _, element := range array {
					element, ok := element.(map[string]interface{})
					if ok {
						_search(element, append(path, key))
					}
				}
			}
		}
	}
	_search(m, make([]string, 0))
	return
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))

	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}
