package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpPost[T ResponseBody[V], V any](baseUrl string, requestBody map[string]interface{}, basicAuthUsernamePassword []string) (_jsonResponseBody T, _statusCode int, _clientError error) {
	var t T
	bb, clientError := json.Marshal(requestBody)
	if clientError != nil {
		return t, 0, clientError
	}

	fullURL := baseUrl
	req, clientError := http.NewRequest("POST", fullURL, bytes.NewBuffer(bb))
	if clientError != nil {
		return t, 0, clientError
	}

	if len(basicAuthUsernamePassword) == 2 {
		req.SetBasicAuth(basicAuthUsernamePassword[0], basicAuthUsernamePassword[1])
	}

	resp, clientError := http.DefaultClient.Do(req)
	if clientError != nil {
		return t, 0, clientError
	}
	defer resp.Body.Close()

	//log.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	body, clientError := io.ReadAll(resp.Body)
	if clientError != nil {
		return t, 0, clientError
	}

	//log.Printf("Body: %s\n", string(body))

	clientError = json.Unmarshal(body, &t)
	if clientError != nil {
		return t, 0, fmt.Errorf("response body: %s %v", string(body), clientError)
	}

	return t, resp.StatusCode, clientError
}

func httpGet[T ResponseBody[V], V any](httpClient *http.Client, fullURL string, queryParams map[string]interface{}) (jsonResponseBody T, statusCode int, _clientError error) {
	var t T

	req, clientError := http.NewRequest("GET", fullURL, nil)
	if clientError != nil {
		return t, 0, clientError
	}

	resp, clientError := httpClient.Do(req)
	if clientError != nil {
		return t, 0, clientError
	}
	defer resp.Body.Close()

	//log.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	body, clientError := io.ReadAll(resp.Body)
	if clientError != nil {
		return t, 0, clientError
	}

	//log.Printf("Body: %s\n", string(body))

	clientError = json.Unmarshal(body, &t)
	if clientError != nil {
		return t, 0, fmt.Errorf("response body: %s %v", string(body), clientError)
	}

	return t, resp.StatusCode, nil
}
