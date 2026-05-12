package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type ChannelService struct {
	fetchClient *http.Client
	ctxTimeout  time.Duration
	repo        *ChannelRepo
	channels    *[]Channel
	mu          sync.RWMutex
	once        sync.Once
}

func NewChannelService(DB *gorm.DB, client *http.Client) *ChannelService {
	repo := NewChannelRepo(DB)
	svc := &ChannelService{
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

func (s *ChannelService) RefreshChannelKey(channelKey *ChannelKey) error {
	payload := channelKey.Metadata.RefreshRequestPayload
	jsonData, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, "POST", channelKey.Metadata.RefreshBaseURL, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	resp, err := s.fetchClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return errors.New("unauthorized err")
		}
		return errors.New("fetch usage err")
	}
	defer func() { _ = resp.Body.Close() }()
	var refresh RefreshResp
	err = json.NewDecoder(resp.Body).Decode(&refresh)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Second * time.Duration(refresh.ExpiresIn)).Unix()
	channelKey.Metadata.RefreshRequestPayload.RefreshToken = refresh.RefreshToken
	channelKey.ApiKey = refresh.AccessToken
	channelKey.Metadata.Exp = exp
	err = s.repo.UpdateChannelKey(channelKey)
	if err != nil {
		return err
	}
	return nil
}
func (s *ChannelService) ColdingChannelKey(channelKey ChannelKey) error {
	return s.repo.UpdateChannelKeyStatus(channelKey.ID, int(Colding))

}
