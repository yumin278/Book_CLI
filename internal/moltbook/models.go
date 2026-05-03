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

// Common Types
type Agent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Karma       int    `json:"karma"`
}

type Submolt struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MemberCount int    `json:"member_count"`
}

type Post struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	URL          string  `json:"url"`
	Submolt      Submolt `json:"submolt"`
	Author       Agent   `json:"author"`
	Upvotes      int     `json:"upvotes"`
	CommentCount int     `json:"comment_count"`
	CreatedAt    string  `json:"created_at"`
}

type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Author    Agent  `json:"author"`
	Upvotes   int    `json:"upvotes"`
	CreatedAt string `json:"created_at"`
}

// Responses

type StatusResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Agent   Agent  `json:"agent"`
}

type FeedResponse struct {
	Success bool   `json:"success"`
	Posts   []Post `json:"posts"`
}

type PostResponse struct {
	Success bool `json:"success"`
	Post    Post `json:"post"`
}

type SubmoltsResponse struct {
	Success  bool      `json:"success"`
	Submolts []Submolt `json:"submolts"`
}

type SearchResult struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Author  Agent  `json:"author,omitempty"`
}

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
}

type DMCheckResponse struct {
	Success     bool   `json:"success"`
	HasActivity bool   `json:"has_activity"`
	Summary     string `json:"summary"`
}

type DMConversation struct {
	ConversationID string `json:"conversation_id"`
	WithAgent      Agent  `json:"with_agent"`
	UnreadCount    int    `json:"unread_count"`
	LastMessageAt  string `json:"last_message_at"`
}

type DMConversationsResponse struct {
	Success       bool `json:"success"`
	Conversations struct {
		Items []DMConversation `json:"items"`
	} `json:"conversations"`
}

type DMMessage struct {
	ID        string `json:"id"`
	From      Agent  `json:"from"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type DMReadResponse struct {
	Success  bool        `json:"success"`
	Messages []DMMessage `json:"messages"`
}

type DMRequest struct {
	ConversationID string `json:"conversation_id"`
	From           Agent  `json:"from"`
	MessagePreview string `json:"message_preview"`
	CreatedAt      string `json:"created_at"`
}

type DMRequestsResponse struct {
	Success  bool        `json:"success"`
	Requests []DMRequest `json:"requests"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
