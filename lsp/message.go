package lsp

import (
	"github.com/google/uuid"
)

type Request struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
	ID     any    `json:"id"` // int32 | string
	Params any    `json:"params,omitempty"`
}

type Response struct {
	RPC   string         `json:"jsonrpc"`
	ID    *any           `json:"id,omitempty"` // int32 | string
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
		ID:     uuid.New().String(),
		Params: params,
	}
}

func NewResponse(id any) Response {
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
