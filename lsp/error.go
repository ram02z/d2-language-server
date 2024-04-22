package lsp

type ErrorCode int

const (
	ParseError                     ErrorCode = -32700
	InvalidRequest                 ErrorCode = -32600
	MethodNotFound                 ErrorCode = -32601
	InvalidParams                  ErrorCode = -32602
	InternalError                  ErrorCode = -32603
	JSONRPCReservedErrorRangeStart ErrorCode = -32099
	ServerNotInitialized           ErrorCode = -32002
	UnknownErrorCode               ErrorCode = -32001
	JSONRPCReservedErrorRangeEnd   ErrorCode = -32000
	LSPReservedErrorRangeStart     ErrorCode = -32899
	RequestFailed                  ErrorCode = -32803
	ServerCancelled                ErrorCode = -32802
	ContentModified                ErrorCode = -32801
	RequestCancelled               ErrorCode = -32800
	LSPReservedErrorRangeEnd       ErrorCode = -32800
)
