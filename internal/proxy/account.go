package proxy

import (
	"sync"
	"time"
)

const (
	Enabled AccountStatus = 1 + iota
	Disabled
	Refreshing
	Colding
)

const (
	APIKeyMode      LoginMode = "api_key"
	AccessTokenMode LoginMode = "access_token"
)

type Account struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string
	TokenExp     int64
	Status       AccountStatus
	UsagePercent float64
	LastCheck    time.Time
	ColdingTime  time.Time
	mu           sync.RWMutex
}

type AccountSnap struct {
	ID           string `json:"id"`
	Name         string
	AccessToken  string
	RefreshToken string
	TokenExp     int64
	Status       AccountStatus
	UsagePercent float64
	LastCheck    time.Time
}

func (a *Account) SnapShot() AccountSnap {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return AccountSnap{
		ID:           a.ID,
		Name:         a.Name,
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		TokenExp:     a.TokenExp,
		Status:       a.Status,
		UsagePercent: a.UsagePercent,
		LastCheck:    a.LastCheck,
	}
}

func (a *Account) UpdateStatus(status AccountStatus) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Status = status
}

func (a *Account) GetStatus() AccountStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Status
}

func (a *Account) GetAccessToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.AccessToken

}

func (a *Account) GetID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.ID

}
func (a *Account) UpdateUsage(usage float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.UsagePercent = usage
}

func (a *Account) UpdateToken(accessToken string, refreshToken string, exp int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.AccessToken = accessToken
	a.RefreshToken = refreshToken
	a.TokenExp = exp
}
