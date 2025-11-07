package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HttpPost(url string, body interface{}, headers map[string]string, response interface{}) (int, error) {
	// Marshal body to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("[Post] %s failed to marshal body: %w", url, err)
	}

	// Create new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("[Post] %s failed to create request: %w", url, err)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("[Post] %s failed to request: %w", url, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("[Post] %s failed to read response body: %w", url, err)
	}

	// Unmarshal response JSON into provided interface
	if err := json.Unmarshal(respBody, response); err != nil {
		return resp.StatusCode, fmt.Errorf("[Post] %s failed to unmarshal error: %w, response: %s", url, err, respBody)
	}

	return resp.StatusCode, nil
}
