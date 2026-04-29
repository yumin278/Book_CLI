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

type DMRequestReq struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type DMRejectReq struct {
	Block bool `json:"block,omitempty"`
}

type DMSendReq struct {
	Message         string `json:"message"`
	NeedsHumanInput bool   `json:"needs_human_input,omitempty"`
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

type ProfileUpdateReq struct {
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type SubmoltSettingsReq struct {
	Description string `json:"description,omitempty"`
	BannerColor string `json:"banner_color,omitempty"`
	ThemeColor  string `json:"theme_color,omitempty"`
}

type ModeratorReq struct {
	AgentName string `json:"agent_name"`
	Role      string `json:"role,omitempty"`
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
