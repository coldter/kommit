package kommit

import (
	"fmt"
	"net/http"
	"strings"
)

type AiClient client

type GenerateGitCommitMsgFromDiffResponse struct {
	Message string `json:"msg"`
	Body    string `json:"body"`
}

func (ac AiClient) GenerateGitCommitMsgFromDiff(diff string) (*GenerateGitCommitMsgFromDiffResponse, error) {
	body := strings.NewReader(diff)

	res, err := ac.client.Post(ac.Url(), body, map[string]string{
		"Content-Type": "text/plain",
	})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseResponseError(res)
	}

	data, err := unmarshalData[*GenerateGitCommitMsgFromDiffResponse](res)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return data, nil
}

func (ac AiClient) Url() string {
	return "/v1/ai.generateGitCommitMsgFromDiff"
}
