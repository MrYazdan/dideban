package types

// PaginationResponse represents pagination metadata in API responses
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Response represents the standard API response wrapper
type Response struct {
	Success    bool                `json:"success"`
	Data       interface{}         `json:"data,omitempty"`
	Error      *Error              `json:"error,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// SuccessResponse creates a successful API response
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// SuccessResponseWithPagination creates a successful API response with pagination
func SuccessResponseWithPagination(data interface{}, pagination *PaginationResponse) Response {
	return Response{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
}
