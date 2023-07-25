package core

import "github.com/lcrownover/hpcadmin-cli/internal/types"

type APIRequest struct {
	Endpoint string
	Method   string
	Body     []byte
}

func NewAPIRequest(endpoint string, method string, data types.APIRequestable) (*APIRequest, error) {
	body, err := data.ToBytes()
	if err != nil {
		return nil, err
	}
	return &APIRequest{
		Endpoint: endpoint,
		Method:   method,
		Body:     body,
	}, nil
}
