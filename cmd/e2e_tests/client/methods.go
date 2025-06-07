package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpPost[T ResponseBody[V], V any](baseUrl string, requestBody map[string]interface{}) (jsonResponseBody T, statusCode int, _clientError error) {
	var t T
	bb, clientError := json.Marshal(requestBody)
	if clientError != nil {
		return t, 0, clientError
	}

	fullURL := baseUrl
	resp, clientError := http.Post(fullURL, "application/json", bytes.NewBuffer(bb))
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
