package schema

type Message struct {
	Role    string
	Content string
}

type CodexResponseRequest struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"messages"`
	CurrentFilePath string    `json:"current_file_path"`
	ProjectName     string    `json:"project_name"`
	RepositoryName  string    `json:"repository_name"`
	CursorLine      uint32    `json:"cursor_line"`
	CursorColum     uint32    `json:"cursor_column"`
	MaxTokens       uint32    `json:"max_tokens"`
	Temperature     float64   `json:"temperature"`
	TopP            uint32    `json:"top_p"`
	N               uint32    `json:"n"`
	Stream          bool      `json:"stream"`
	Stop            []string  `json:"stop"`
}

type CodexResponseStreamChunk struct {
	ID string `json:"id,omitempty"`
}

type UsageResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	PlanType string `json:"plan_type"` // "free", "pro", "plus"

	// 核心配额
	RateLimit struct {
		Allowed       bool `json:"allowed"`
		LimitReached  bool `json:"limit_reached"`
		PrimaryWindow struct {
			UsedPercent        float64 `json:"used_percent"`         // 已使用百分比
			LimitWindowSeconds int64   `json:"limit_window_seconds"` // 总窗口（如 7天）
			ResetAfterSeconds  int64   `json:"reset_after_seconds"`  // 剩余冷却秒数
			ResetAt            int64   `json:"reset_at"`             // 重置时间戳
		} `json:"primary_window"`
	} `json:"rate_limit"`

	// 代码审查配额（通常与上面分开）
	CodeReviewRateLimit struct {
		Allowed       bool `json:"allowed"`
		LimitReached  bool `json:"limit_reached"`
		PrimaryWindow struct {
			UsedPercent float64 `json:"used_percent"`
			ResetAt     int64   `json:"reset_at"`
		} `json:"primary_window"`
	} `json:"code_review_rate_limit"`
}
