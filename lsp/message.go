package lsp

import (
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

type RequestID string

func (id *RequestID) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*id = RequestID(s)
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*id = RequestID(strconv.Itoa(i))
	return nil
}

type Request struct {
	RPC    string    `json:"jsonrpc"`
	Method string    `json:"method"`
	ID     RequestID `json:"id"`
	Params any       `json:"params,omitempty"`
}

type Response struct {
	RPC   string         `json:"jsonrpc"`
	ID    *RequestID     `json:"id,omitempty"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}

func NewRequest(method Method) Request {
	return NewRequestWithParams(method, nil)
}

func NewRequestWithParams(method Method, params any) Request {
	return Request{
		RPC:    JsonRpc,
		Method: string(method),
		ID:     RequestID(uuid.New().String()),
		Params: params,
	}
}

func NewResponse(id RequestID) Response {
	return Response{
		RPC: JsonRpc,
		ID:  &id,
	}
}

func NewNotification(method Method) Notification {
	return Notification{
		RPC:    JsonRpc,
		Method: string(method),
	}
}
