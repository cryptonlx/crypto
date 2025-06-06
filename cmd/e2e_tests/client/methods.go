package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func httpPost[T ResponseBody[V], V any](baseUrl string, requestBody map[string]interface{}) (jsonResponseBody T, statusCode int, clientError error) {
	var t T
	bb, err := json.Marshal(requestBody)
	if err != nil {
		return t, 0, err
	}

	fullURL := baseUrl
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(bb))
	if err != nil {
		return t, 0, err
	}
	defer resp.Body.Close()

	//log.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t, 0, err
	}

	//log.Printf("Body: %s\n", string(body))

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, 0, fmt.Errorf("response body: %s %v", string(body), err)
	}

	return t, resp.StatusCode, err
}
