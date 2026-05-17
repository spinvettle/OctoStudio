package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	ErrChannelKeyRefreshInProgress = errors.New("channel key refresh already in progress")
	ErrRefreshTokenUnauthorized    = errors.New("refresh token unauthorized")
	ErrRefreshTokenFetchFailed     = errors.New("refresh token fetch failed")
)

type ChannelService struct {
	fetchClient *http.Client
	Cfg         ChannelServiceConfig
	repo        *ChannelRepo
	channels    *[]Channel
	mu          sync.RWMutex
	once        sync.Once
	group       singleflight.Group
}

type ChannelServiceConfig struct {
	FetchTimeout time.Duration
}

func NewChannelService(DB *gorm.DB, client *http.Client, cfg *ChannelServiceConfig) *ChannelService {
	repo := NewChannelRepo(DB)
	svc := &ChannelService{
		Cfg:         *cfg,
		repo:        repo,
		fetchClient: client,
	}
	svc.autoRefreshChannelsCache()
	return svc
}

func (s *ChannelService) GetAllChannel() (*[]Channel, error) {
	allChannels, err := s.repo.GetAllChannel()
	if err != nil {
		return nil, err
	}
	return allChannels, nil
}
func (s *ChannelService) autoRefreshChannelsCache() {
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for range ticker.C {
			s.autoRefreshChannelsCache()
		}
	}()
}
func (s *ChannelService) LoadChannelsCache() {
	var err error
	s.mu.Lock()
	defer s.mu.Unlock()
	s.channels, err = s.GetAllChannel()
	if err != nil {
		slog.Error("failed to init channels cache", slog.Any("error", err))
	}
}

func (s *ChannelService) RandomPickChannelByModel(modelName string) *[]Channel {

	s.once.Do(s.LoadChannelsCache)

	var selectedChannesl []Channel
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range *s.channels {
		for _, m := range ch.Models {
			if strings.ToLower(modelName) == m.Name {
				selectedChannesl = append(selectedChannesl, ch)
			}
		}
	}
	return &selectedChannesl

}
func (s *ChannelService) fetchNewToken(ctx context.Context, path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := s.fetchClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ChannelService) RefreshKeyOnce(channelKey *ChannelKey) (bool, error) {
	key := strconv.Itoa(channelKey.ID)
	_, err, shared := s.group.Do(key, func() (any, error) { // singleflight do refreshing
		err := s.refreshChannelKey(channelKey)
		return nil, err
	})

	return shared, err

}
func (s *ChannelService) refreshChannelKey(channelKey *ChannelKey) error {
	var refreshingSuccess bool
	success, err := s.repo.UpdateChannelKeyStatusWithCondition(channelKey.ID, int(Enable), int(Refreshing))
	if err != nil {
		return err
	}
	if !success { //刷新状态设置失败，已经有操作在执行
		return ErrChannelKeyRefreshInProgress
	}
	defer func() { //如果success设置Enable->Refreshing,现在失败了需要设置成Disable
		if !refreshingSuccess {
			_, _ = s.repo.UpdateChannelKeyStatusWithCondition(channelKey.ID, int(Refreshing), int(Disable))
		}
	}()
	payload := channelKey.Metadata.RefreshRequestPayload
	jsonData, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.FetchTimeout)
	defer cancel()
	resp, err := s.fetchNewToken(ctx, channelKey.Metadata.RefreshBaseURL, jsonData)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrRefreshTokenUnauthorized
		}
		return fmt.Errorf("%w: status %d", ErrRefreshTokenFetchFailed, resp.StatusCode)
	}

	var refresh RefreshResp
	err = json.NewDecoder(resp.Body).Decode(&refresh)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Second * time.Duration(refresh.ExpiresIn)).Unix()
	channelKey.Metadata.RefreshRequestPayload.RefreshToken = refresh.RefreshToken
	channelKey.ApiKey = refresh.AccessToken
	channelKey.Metadata.Exp = exp
	channelKey.Status = Enable
	err = s.repo.UpdateChannelKey(channelKey)
	if err != nil {
		return err
	}
	refreshingSuccess = true
	return nil
}
func (s *ChannelService) ColdingChannelKey(channelKey ChannelKey) error {
	err := s.repo.UpdateChannelKeyStatus(channelKey.ID, int(Colding))
	return err

}

func (s *ChannelService) UpdateChannelKeyStatus(id int, status ChannelKeyStatus) error {
	return s.repo.UpdateChannelKeyStatus(id, int(status))
}
