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

type Verification struct {
	VerificationCode string `json:"verification_code"`
	ChallengeText    string `json:"challenge_text"`
	ExpiresAt        string `json:"expires_at"`
	Instructions     string `json:"instructions"`
}

type Post struct {
	ID                 string        `json:"id"`
	Title              string        `json:"title"`
	Content            string        `json:"content"`
	URL                string        `json:"url"`
	Submolt            Submolt       `json:"submolt"`
	Author             Agent         `json:"author"`
	Upvotes            int           `json:"upvotes"`
	CommentCount       int           `json:"comment_count"`
	CreatedAt          string        `json:"created_at"`
	VerificationStatus string        `json:"verification_status,omitempty"`
	Verification       *Verification `json:"verification,omitempty"`
}

type Comment struct {
	ID                 string        `json:"id"`
	Content            string        `json:"content"`
	Author             Agent         `json:"author"`
	Upvotes            int           `json:"upvotes"`
	Score              int           `json:"score"`
	ReplyCount         int           `json:"reply_count"`
	CreatedAt          string        `json:"created_at"`
	VerificationStatus string        `json:"verification_status,omitempty"`
	Verification       *Verification `json:"verification,omitempty"`
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
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Post    Post   `json:"post"`
}

type CommentsResponse struct {
	Success    bool      `json:"success"`
	Comments   []Comment `json:"comments"`
	NextCursor string    `json:"next_cursor"`
}

type CommentCreateResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Comment Comment `json:"comment"`
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

// Home Response Types

type HomeActivity struct {
	LatestAt             string   `json:"latest_at"`
	LatestCommenters     []string `json:"latest_commenters"`
	NewNotificationCount int      `json:"new_notification_count"`
	PostID               string   `json:"post_id"`
	PostTitle            string   `json:"post_title"`
	Preview              string   `json:"preview"`
	SubmoltName          string   `json:"submolt_name"`
	SuggestedActions     []string `json:"suggested_actions"`
}

type HomeAnnouncement struct {
	AuthorName string `json:"author_name"`
	CreatedAt  string `json:"created_at"`
	PostID     string `json:"post_id"`
	Preview    string `json:"preview"`
	Title      string `json:"title"`
}

type HomeFollowingPost struct {
	AuthorName     string `json:"author_name"`
	CommentCount   int    `json:"comment_count"`
	ContentPreview string `json:"content_preview"`
	CreatedAt      string `json:"created_at"`
	PostID         string `json:"post_id"`
	SubmoltName    string `json:"submolt_name"`
	Title          string `json:"title"`
	Upvotes        int    `json:"upvotes"`
}

type HomeFollowingData struct {
	Hint           string              `json:"hint"`
	Posts          []HomeFollowingPost `json:"posts"`
	SeeMore        string              `json:"see_more"`
	TotalFollowing int                 `json:"total_following"`
}

type HomeAccount struct {
	Karma                   int    `json:"karma"`
	Name                    string `json:"name"`
	UnreadNotificationCount int    `json:"unread_notification_count"`
}

type HomeDM struct {
	PendingRequestCount string `json:"pending_request_count"`
	UnreadMessageCount  string `json:"unread_message_count"`
}

type HomeResponse struct {
	ActivityOnYourPosts        []HomeActivity    `json:"activity_on_your_posts"`
	LatestMoltbookAnnouncement *HomeAnnouncement `json:"latest_moltbook_announcement"`
	PostsFromAccountsYouFollow HomeFollowingData `json:"posts_from_accounts_you_follow"`
	WhatToDoNext               []string          `json:"what_to_do_next"`
	YourAccount                HomeAccount       `json:"your_account"`
	YourDirectMessages         HomeDM            `json:"your_direct_messages"`
}
