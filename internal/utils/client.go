package utils

import (
	"github.com/coldter/kommit/kommit"
	"github.com/coldter/kommit/settings"
	"net/url"
)

const (
	//kommitRpcDefaultBaseUrl = "http://localhost:3100"
	kommitRpcDefaultBaseUrl = "https://kommit-ai-rpc.kuldeep.workers.dev"
)

func GetClient() (*kommit.Client, error) {
	s, err := settings.ReadSettings()
	if err != nil {
		return nil, err
	}

	configUrl := s.GetBaseURL()
	if configUrl == "" {
		configUrl = kommitRpcDefaultBaseUrl
	}
	u, err := url.Parse(configUrl)
	return kommit.New(u, s.GetToken()), nil
}
