package moltbook

// PostCreateReq represents a request to create a post
// submolt_name is the canonical API field; server may also accept submolt alias.
type PostCreateReq struct {
	SubmoltName string `json:"submolt_name"`
	Title       string `json:"title"`
	Content     string `json:"content,omitempty"`
	URL         string `json:"url,omitempty"`
}

// CommentCreateReq represents a request to create a comment
type CommentCreateReq struct {
	Content  string `json:"content"`
	ParentID string `json:"parent_id,omitempty"`
}

// VerifyReq submits a verification challenge answer
type VerifyReq struct {
	VerificationCode string `json:"verification_code"`
	Answer           string `json:"answer"`
}

// RegisterReq represents a request to register a new agent
type RegisterReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Success        bool   `json:"success"`
	Error          string `json:"error"`
	Hint           string `json:"hint,omitempty"`
	RetryAfter     int    `json:"retry_after,omitempty"`
	RetryAfterSecs int    `json:"retry_after_seconds,omitempty"`
	RetryAfterMins int    `json:"retry_after_minutes,omitempty"`
	DailyRemaining int    `json:"daily_remaining,omitempty"`
}
