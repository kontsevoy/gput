/*
Rackspace API wrapper which includes:
	- Authentication
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	identityUrl = "https://identity.api.rackspacecloud.com/v2.0/tokens"
	authRequest = `{ "auth": {
				"RAX-KSKEY:apiKeyCredentials": {
					"apiKey": "%s",
					"username": "kontsevoy.api"
				}}}`
)

// This data structure is returned as JSON when we authenticate.
// Generated by ChimeraCoder/gojson from curl/auth_output.json
type RaxSession struct {
	Access struct {
		ServiceCatalog []struct {
			Endpoints []struct {
				PublicURL string `json:"publicURL"`
				Region    string `json:"region"`
				TenantID  string `json:"tenantId"`
			} `json:"endpoints"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"serviceCatalog"`
		Token struct {
			RAX_AUTH_authenticatedBy []string `json:"RAX-AUTH:authenticatedBy"`
			Expires                  string   `json:"expires"`
			ID                       string   `json:"id"`
			Tenant                   struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"tenant"`
		} `json:"token"`
		User struct {
			RAX_AUTH_defaultRegion string `json:"RAX-AUTH:defaultRegion"`
			ID                     string `json:"id"`
			Name                   string `json:"name"`
			Roles                  []struct {
				Description string `json:"description"`
				ID          string `json:"id"`
				Name        string `json:"name"`
				TenantID    string `json:"tenantId"`
			} `json:"roles"`
		} `json:"user"`
	} `json:"access"`
}

// Authenticates against Rackspace Auth and returns an authentication token
func authenticate(apiKey string) (session RaxSession, err error) {
	// HTTP POST to auth URL:
	requestBody := fmt.Sprintf(authRequest, apiKey)
	response, err := http.Post(identityUrl, "application/json",
		strings.NewReader(requestBody))
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(responseBody, &session)
	if err != nil {
		return
	}
	return
}

// Extracts an entry point for a given region + service type from the JSON
// data returned by Auth API call
func (ai *RaxSession) getEntryPoint(region string, serviceType string) (entryPoint string) {
	for _, service := range ai.Access.ServiceCatalog {
		if service.Type == serviceType {
			for _, ep := range service.Endpoints {
				if ep.Region == region {
					entryPoint = ep.PublicURL
					return entryPoint
				}
			}
		}
	}
	return
}

func (ai *RaxSession) listContainers(url string) (err error) {
	request, err := http.NewRequest("GET", url+"?format=json", nil)
	if err != nil {
		return
	}
	request.Header.Add("X-Auth-Token", ai.Access.Token.ID)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("--> %v <--\n", string(body))

	request, err = http.NewRequest("GET", url+"/ev-public", nil)
	panicIf(err)
	request.Header.Add("X-Auth-Token", ai.Access.Token.ID)

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	body, _ = ioutil.ReadAll(response.Body)
	fmt.Printf("--> %v <--\n", string(body))
	return
}

// Create/update CloudFiles object (file)
func (ai *RaxSession) upsertObject(url string, r io.Reader, objectName string) {
	url = strings.Join([]string{url, "ev-private", objectName}, "/")
	fmt.Println(url)

	request, err := http.NewRequest("PUT", url, r)
	panicIf(err)

	request.Header.Add("X-Auth-Token", ai.Access.Token.ID)
	request.Header.Add("Content-Type", "text/plain")

	fmt.Println(request.Header)

	response, err := http.DefaultClient.Do(request)
	panicIf(err)

	fmt.Println(response)
}