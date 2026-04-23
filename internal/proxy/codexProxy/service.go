package codexProxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"maps"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

var (
	ErrNotFoundAvailabelAccount = errors.New("not found any available account")
	ErrEmptyPool                = errors.New("account pool is empty")
	ErrUnauthorized             = errors.New("unauthorized err")
	ErrAccountNotExists         = errors.New("account not exists")
	ErrRelayDefaut              = errors.New("relay error")
)

const CodexResponsesURL = "https://chatgpt.com/backend-api/codex/responses"
const CodexUsageURL = "https://chatgpt.com/backend-api//wham/usage"
const CodexRefreshTokenURL = "https://auth.openai.com/oauth/token"
const baseBackoff = time.Millisecond * 500
const maxBackoff = time.Second * 10

var ProxyService *codexProxyService

func InitCodexProxy() {
	var err error
	InitHttpClient()
	config := &CodexProxyConfig{
		FetchClient:  FetchClient,
		RelayClient:  RelayClient,
		AccountsFile: "./accounts.json",
	}
	ProxyService, err = NewcodexProxyService(config)
	if err != nil {
		panic(err)
	}
	_, _ = ProxyService.AddBatchAccounts(ProxyService.config.AccountsFile)
}

type codexProxyService struct {
	pool         *AccountPool
	openaiClient *openaiClient
	relayClient  *http.Client
	config       *CodexProxyConfig
}

type CodexProxyConfig struct {
	MaxRetry           int64
	AccountsFile       string
	ColdingTime        time.Duration
	OpenaifetchTimeOut time.Duration
	CodexRelayTimeOut  time.Duration
	FetchClient        *http.Client
	RelayClient        *http.Client
}

func NewcodexProxyService(config *CodexProxyConfig) (*codexProxyService, error) {
	if config.ColdingTime == 0 {
		config.ColdingTime = time.Second * 5
	}
	if config.OpenaifetchTimeOut == 0 {
		config.OpenaifetchTimeOut = time.Second * 60
	}
	if config.CodexRelayTimeOut == 0 {
		config.CodexRelayTimeOut = time.Minute * 10
	}
	if config.MaxRetry == 0 {
		config.MaxRetry = 5
	}
	if config.FetchClient == nil || config.RelayClient == nil {
		return nil, errors.New("nil *http.Client is not allowed")
	}
	return &codexProxyService{
		pool:         NewAccountPool(),
		openaiClient: NewOpenaiClient(config.FetchClient),
		relayClient:  config.RelayClient,
		config:       config,
	}, nil
}

func (s *codexProxyService) AddBatchAccounts(path string) (int, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var accountFile AccountJsonFile
	err = json.Unmarshal(bytes, &accountFile)
	if err != nil {
		return 0, nil
	}
	if accountFile.Type != "Codex" {
		return 0, errors.New("Unspported type:" + accountFile.Type)
	}
	var wg sync.WaitGroup
	var count int
	for i := 0; i < len(accountFile.Accounts); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			accountInfo := accountFile.Accounts[i]
			err := s.AddAccount(accountInfo.Name, "", accountInfo.AccessToken, accountInfo.RefreshToken)
			if err != nil {
				slog.Error("failed to add account",
					slog.String("name", accountInfo.Name),
					slog.Any("error", err))
				return
			}
			count += 1

		}(i)
	}
	wg.Wait()
	return count, nil
}

func (s *codexProxyService) saveToDisk(ctx context.Context) {
	//TODO 数据库
	accounts := s.pool.GetAllAccounts()
	var items []AccountItem
	for _, item := range *accounts {
		items = append(items, AccountItem{
			Name:         item.Name,
			AccessToken:  item.AccessToken,
			RefreshToken: item.RefreshToken,
		})
	}
	accountFile := AccountJsonFile{
		Type:     "Codex",
		Accounts: items,
	}
	bytes, err := json.MarshalIndent(accountFile, "", " ")
	if err != nil {
		slog.ErrorContext(ctx, "marshall json failed", slog.Any("error", err))
	}
	err = os.WriteFile(s.config.AccountsFile, bytes, 0644)
	if err != nil {
		slog.ErrorContext(ctx, "save account failed", slog.Any("error", err))
	}
}

func (s *codexProxyService) buildRequest(ctx context.Context, headers map[string][]string, body []byte, accessToken string) (*http.Request, error) {
	codexProxyReq, err := http.NewRequestWithContext(ctx, "POST", CodexResponsesURL, bytes.NewBuffer(body))
	if err != nil {

		return nil, err
	}
	maps.Copy(codexProxyReq.Header, headers)
	codexProxyReq.Header["Authorization"] = []string{"Bearer " + accessToken}
	return codexProxyReq, nil
}

func (s *codexProxyService) DoProxyRequest(ctx context.Context, body []byte, headers map[string][]string) (*http.Response, error) {
	var resp *http.Response
	var account *Account
	var err error
	needAccount := true
	for attempt := 0; attempt < int(s.config.MaxRetry); attempt++ {
		if needAccount {
			account, err = s.pool.GetAccount()
			if err != nil {
				return nil, err
			}
		}
		AccountSnap := account.SnapShot()

		accessToken := AccountSnap.AccessToken
		ID := AccountSnap.ID
		slog.InfoContext(ctx, "choice account", slog.String("id", AccountSnap.ID), slog.String("name", AccountSnap.Name))

		codexProxyReq, err := s.buildRequest(ctx, headers, body, accessToken)
		if err != nil {
			return nil, err
		}
		resp, err = s.relayClient.Do(codexProxyReq)
		if err != nil {
			slog.ErrorContext(ctx, "relay client do failed", slog.Any("error", err))
			return nil, err
		}
		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
		case http.StatusUnauthorized, http.StatusForbidden:
			//刷新账号
			account.UpdateStatus(Refreshing) //先手动设置一遍，防止异步更新状态不及时
			needAccount = true
			slog.InfoContext(ctx, "need to fresh account", slog.String("id", AccountSnap.ID), slog.String("name", AccountSnap.Name))
			ctx, cancel := context.WithTimeout(ctx, s.config.OpenaifetchTimeOut)
			go func() {
				defer cancel()
				s.Refresh(ctx, ID)
			}()

		case http.StatusRequestTimeout,
			http.StatusGatewayTimeout,
			http.StatusServiceUnavailable:

			slog.InfoContext(ctx, "relay need to wait", slog.String("reason", resp.Status), slog.Int("retry attempt", attempt))
			needAccount = false
			sleepTime := baseBackoff * (1 << attempt)

			sleepTime += time.Duration(rand.Int63n(int64(baseBackoff) / 5))
			if sleepTime > maxBackoff {
				sleepTime = maxBackoff
			}
			time.Sleep(sleepTime)

		case http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusUnprocessableEntity:
			slog.InfoContext(ctx, "upstream returned client error", slog.String("reason", resp.Status))
			return resp, nil

		case http.StatusTooManyRequests:
			//冷却账号
			slog.InfoContext(ctx, "need to cold account", slog.String("id", AccountSnap.ID), slog.String("name", AccountSnap.Name))
			account.UpdateStatus(Colding)
			needAccount = true
			go s.Colding(ID)

		default:
			return resp, ErrRelayDefaut

		}

		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				slog.ErrorContext(ctx, "failed to close upstream response body", slog.Any("error", err))
			}
		}

	}
	return resp, ErrRelayDefaut
}

func (s *codexProxyService) AddAccount(name, apiKey, accessToken, refreshToken string) error {
	exp, _, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return err
	}
	var account = &Account{
		Name:         name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenExp:     exp,
	}

	account.ID = uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	usage, err := s.openaiClient.fetchUsage(ctx, account.AccessToken)
	if err != nil {
		return err
	}
	account.UsagePercent = 100 - usage
	err = s.pool.AddAccount(account)
	if err != nil {
		return err
	}

	return nil

}

func (s *codexProxyService) UpdateAccount(id, apiKey, accessToken, refreshToken string) error {
	return nil
}

func (s *codexProxyService) DeleteAccount(id string) error {
	return nil
}

func (s *codexProxyService) Refresh(ctx context.Context, id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get account by id", slog.String("id", id), slog.Any("error", err))
		return
	}
	accountSnap := account.SnapShot()
	account.UpdateStatus(Refreshing)
	accessToken, refreshToken, exp, err := s.openaiClient.fetchToken(ctx, accountSnap.RefreshToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch new access token",
			slog.String("id", id),
			slog.String("name", accountSnap.Name),
			slog.Any("error", err))
		account.UpdateStatus(Disabled)
		return
	}
	account.UpdateToken(accessToken, refreshToken, exp)
	account.UpdateStatus(Enabled)

	s.saveToDisk(ctx)

}

func (s *codexProxyService) Colding(id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		slog.Error("failed to get account for colding", slog.String("id", id), slog.Any("error", err))
		return
	}
	account.UpdateStatus(Colding)
	// accountSnap := account.SnapShot()
	time.AfterFunc(s.config.ColdingTime, func() {
		account.UpdateStatus(Enabled)
	})

}

func (s *codexProxyService) GetUsage(id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		slog.Error("failed to get account for usage", slog.String("id", id), slog.Any("error", err))
		return
	}
	accountSnap := account.SnapShot()
	ctx, cancel := context.WithTimeout(context.Background(), s.config.OpenaifetchTimeOut)
	defer cancel()
	usage, err := s.openaiClient.fetchUsage(ctx, accountSnap.AccessToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch usage",
			slog.String("id", id),
			slog.String("name", accountSnap.Name),
			slog.Any("error", err))
		return
	}
	account.UpdateUsage(usage)

}

type openaiClient struct {
	httpClient *http.Client
}

func NewOpenaiClient(client *http.Client) *openaiClient {
	return &openaiClient{httpClient: client}
}
func (c *openaiClient) fetchToken(ctx context.Context, refreshToken string) (string, string, int64, error) {

	var req = RefreshReq{
		IdClient:     "app_EMoamEEZ73f0CkXaXp7hrann",
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	jsonData, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, "POST", CodexRefreshTokenURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", "", 0, err
	}
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", "", 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", "", 0, ErrUnauthorized
		}
		return "", "", 0, errors.New("fetch usage err")
	}
	var refresh RefreshResp
	err = json.NewDecoder(resp.Body).Decode(&refresh)
	if err != nil {
		return "", "", 0, err
	}

	exp := time.Now().Add(time.Second * time.Duration(refresh.ExpiresIn))
	//ExpiresIn返回的是秒，转成unix

	return refresh.AccessToken, refresh.RefreshToken, int64(exp.Unix()), nil

}

func (c *openaiClient) fetchUsage(ctx context.Context, accessToken string) (float64, error) {
	var usageResp Usage
	request, _ := http.NewRequestWithContext(ctx, "GET", CodexUsageURL, nil)
	request.Header.Set("Host", "chatgpt.com")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return 0, ErrUnauthorized
		}
		return 0, errors.New("fetch usage err")
	}
	err = json.NewDecoder(resp.Body).Decode(&usageResp)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode usage response", slog.Any("error", err))
		return 0, err
	}
	return usageResp.RateLimit.PrimaryWindow.UsedPercent, nil
}
