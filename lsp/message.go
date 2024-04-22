package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
	ID     int    `json:"id"`
}

type Response struct {
	RPC   string         `json:"jsonrpc"`
	ID    *int           `json:"id,omitempty"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}

func NewResponse(id int) Response {
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
