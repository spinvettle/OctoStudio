package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"maps"
	"net/http"
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

var proxyService *ProxyService

type ProxyService struct {
	pool               *AccountPool
	openaiClient       *openaiClient
	relayClient        *http.Client
	accountColdingTime time.Duration
}

func NewProxyService(relayClient *http.Client, openaiClient *http.Client, coldingTime time.Duration) *ProxyService {
	return &ProxyService{
		pool:               NewAccountPool(),
		openaiClient:       NewOpenaiClient(openaiClient),
		relayClient:        relayClient,
		accountColdingTime: coldingTime,
	}
}

func buildRequest(headers map[string][]string, body []byte, accessToken string) (*http.Request, error) {
	proxyReq, err := http.NewRequest("POST", CodexResponsesURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	maps.Copy(proxyReq.Header, headers)
	proxyReq.Header["Authorization"] = []string{"Bearer " + accessToken}
	return proxyReq, nil
}

func (s *ProxyService) DoProxyRequest(body []byte, headers map[string][]string) (*http.Response, error) {
	var resp *http.Response
	var account *Account
	var err error
	// firstTime:=true
	needAccount := true
	//暂时写死5次
	for range 5 {
		if needAccount {
			account, err = s.pool.GetAccount()
			if err != nil {
				return nil, err
			}
		}

		accessToken := account.GetAccessToken()
		ID := account.GetID()

		proxyReq, err := buildRequest(headers, body, accessToken)
		if err != nil {
			return nil, err
		}

		resp, err = s.relayClient.Do(proxyReq)
		if err != nil {
			return nil, err
		}
		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
		case http.StatusUnauthorized, http.StatusForbidden:
			//刷新账号
			account.UpdateStatus(Refreshing) //先手动设置一遍，防止异步更新状态不及时
			needAccount = true
			go s.Refresh(ID)

		case http.StatusRequestTimeout,
			http.StatusGatewayTimeout,
			http.StatusServiceUnavailable:
			//TDO退避重试
			needAccount = false
			time.Sleep(time.Millisecond * 100)

		case http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusUnprocessableEntity:

			return resp, nil

		case http.StatusTooManyRequests:
			//冷却账号
			account.UpdateStatus(Colding)
			needAccount = true
			go s.Colding(ID)

		default:
			return resp, ErrRelayDefaut

		}

		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Printf("close body error: %v", err)
			}
		}

	}
	return resp, ErrRelayDefaut
}

func (s *ProxyService) AddAccount(name, apiKey, accessToken, refreshToken string) error {
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
	account.UsagePercent = usage
	err = s.pool.AddAccount(account)
	if err != nil {
		return err
	}

	return nil

}

func (s *ProxyService) UpdateAccount(id, apiKey, accessToken, refreshToken string) error {
	return nil
}

func (s *ProxyService) DeleteAccount(id string) error {
	return nil
}

func (s *ProxyService) Refresh(id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		log.Println("Refresh err:" + err.Error())
		return
	}
	account.UpdateStatus(Refreshing)
	accountSnap := account.SnapShot()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	accessToken, refreshToken, exp, err := s.openaiClient.fetchToken(ctx, accountSnap.RefreshToken)
	if err != nil {
		log.Println("Refresh err:" + err.Error())
		account.UpdateStatus(Disabled)
		return
	}
	account.UpdateToken(accessToken, refreshToken, exp)
	account.UpdateStatus(Enabled)

}

func (s *ProxyService) Colding(id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		log.Println("Refresh err:" + err.Error())
		return
	}
	account.UpdateStatus(Colding)
	// accountSnap := account.SnapShot()
	time.AfterFunc(s.accountColdingTime, func() {
		account.UpdateStatus(Enabled)
	})

}

func (s *ProxyService) GetUsage(id string) {
	account, err := s.pool.GetAccountById(id)
	if err != nil {
		log.Println("Refresh err:" + err.Error())
	}
	accountSnap := account.SnapShot()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	usage, err := s.openaiClient.fetchUsage(ctx, accountSnap.AccessToken)
	if err != nil {
		log.Println("Refresh err:" + err.Error())
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
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
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
		log.Printf("请求解析错误: %v", err)
		return 0, err
	}
	return usageResp.RateLimit.PrimaryWindow.UsedPercent, nil
}
