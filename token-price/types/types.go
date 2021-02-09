package types

const (
	ServiceName = "token-price"
)

type RequestCallback func(reqID, input string) (output *ServiceOutput, requestResult *RequestResult)

type State int

// Status returned of RequestCallback
const (
	Success = iota
	ClientError
	ServiceError
)

// RequestResult is result of RequestCallback
type RequestResult struct {
	State   State // Use status returned
	Message string
}

// Result of RequestCallback
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}
