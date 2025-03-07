package kommit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func unmarshalAPIResponse[T any](r *http.Response) (APIResponse[T], error) {
	var response APIResponse[T]

	if r.Body == nil {
		return response, fmt.Errorf("empty response body")
	}
	defer r.Body.Close()

	// Read the entire response body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if we have any content
	if len(bodyBytes) == 0 {
		return response, fmt.Errorf("empty response body")
	}

	// Unmarshal directly into our APIResponse struct
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return response, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

func unmarshalData[T any](r *http.Response) (T, error) {
	var empty T

	apiResponse, err := unmarshalAPIResponse[T](r)
	if err != nil {
		return empty, err
	}

	return apiResponse.Data, nil
}

func parseResponseError(res *http.Response) error {
	type ErrorResponse struct {
		Error   interface{} `json:"error,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	e := new(ErrorResponse)
	err := json.NewDecoder(res.Body).Decode(e)
	if err != nil {
		return err
	}

	return fmt.Errorf("error: %v, message: %s", e.Error, e.Message)
}
