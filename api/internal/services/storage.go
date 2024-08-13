package services

import (
	"io"
	"net/http"
)

type StorageClient struct {
	endpoint   string
	httpClient *http.Client
}

func (s *StorageClient) GetObject(key string) ([]byte, error) {
	path := s.endpoint + "/" + key
	resp, err := s.httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func NewStorageClient(endpoint string) *StorageClient {
	return &StorageClient{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}
}
