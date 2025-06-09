package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var postLock sync.Mutex

func httpPost[T ResponseBody[V], V any](httpClient *http.Client, baseUrl string, requestBody map[string]interface{}, basicAuthUsernamePassword []string) (_jsonResponseBody T, _statusCode int, _clientError error) {
	postLock.Lock()
	time.Sleep(1 * time.Millisecond)
	defer postLock.Unlock()
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

	return t, resp.StatusCode, clientError
}

var getLock sync.Mutex

func httpGet[T ResponseBody[V], V any](httpClient *http.Client, fullURL string, queryParams map[string]interface{}) (jsonResponseBody T, statusCode int, _clientError error) {
	getLock.Lock()
	time.Sleep(1 * time.Millisecond)
	defer getLock.Unlock()
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
