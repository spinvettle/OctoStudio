package proxy

type RateLimitWindow struct {
	UsedPercent        float64 `json:"used_percent"`
	LimitWindowSeconds int     `json:"limit_window_seconds"`
	ResetAfterSeconds  int     `json:"reset_after_seconds"`
	ResetAt            int64   `json:"reset_at"`
}

type RateLimitInfo struct {
	Allowed         bool            `json:"allowed"`
	LimitReached    bool            `json:"limit_reached"`
	PrimaryWindow   RateLimitWindow `json:"primary_window"`
	SecondaryWindow interface{}     `json:"secondary_window"` // 因为 JSON 里是 null，用 interface{} 接收
}

type Usage struct {
	UserID               string         `json:"user_id"`
	AccountID            string         `json:"account_id"`
	Email                string         `json:"email"`
	PlanType             string         `json:"plan_type"`
	RateLimit            RateLimitInfo  `json:"rate_limit"`
	CodeReviewRateLimit  RateLimitInfo  `json:"code_review_rate_limit"`
	AdditionalRateLimits interface{}    `json:"additional_rate_limits"`
	Credits              map[string]any `json:"credits"` // 也可以用 map 或自定义结构体
	SpendControl         struct {
		Reached bool `json:"reached"`
	} `json:"spend_control"`
}

type RefreshReq struct {
	IdClient     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}
