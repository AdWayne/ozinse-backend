package model

type APIError struct {
	ErrorCode string                 `json:"error_code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
}

func NewAPIError(code, message string, details map[string]interface{}) APIError {
	if details == nil {
		details = make(map[string]interface{})
	}
	return APIError{
		ErrorCode: code,
		Message:   message,
		Details:   details,
	}
}