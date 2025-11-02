package dto

// Response is a generic API response wrapper
type Response struct {
	Data interface{} `json:"data,omitempty"`
	Meta *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination and other metadata
type Meta struct {
	Page   int `json:"page"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewResponse creates a new response with data
func NewResponse(data interface{}) Response {
	return Response{
		Data: data,
	}
}

// NewListResponse creates a new response with data and metadata
func NewListResponse(data interface{}, page, limit, offset int) Response {
	return Response{
		Data: data,
		Meta: &Meta{
			Page:   page,
			Limit:  limit,
			Offset: offset,
		},
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
