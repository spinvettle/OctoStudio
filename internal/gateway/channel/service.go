package channel

import (
	"log/slog"
	"strings"
	"sync"
	"time"
)

type ChannelService struct {
	// client   *http.Client
	repo     *ChannelRepo
	channels *[]Channel
	mu       sync.RWMutex
	once     sync.Once
}

func NewChannelService(repo *ChannelRepo) *ChannelService {
	svc := &ChannelService{repo: repo}
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

func (s *ChannelService) RefreshChannelKey(channelKey ChannelKey) error {
	return nil
}
func (s *ChannelService) ColdingChannelKey(channelKey ChannelKey) error {
	channelKey.Status = Colding
	return s.repo.UpdateChannelKey(&channelKey)

}
