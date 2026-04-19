package codexProxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spinvettle/OctoStudio/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
func setUpMockHttpClient(sleepMil time.Duration, fn RoundTripFunc) *http.Client {
	time.Sleep(sleepMil)
	return &http.Client{
		Transport: fn,
	}

}

func usageStatusOK(req *http.Request) *http.Response {
	usage := Usage{RateLimit: RateLimitInfo{
		PrimaryWindow: RateLimitWindow{
			UsedPercent: 10.0,
		},
	}}
	bytesData, _ := json.Marshal(usage)
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(bytesData)),
	}
}

func StatusUnauthorized(req *http.Request) *http.Response {

	return &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     make(http.Header),
		Body: io.NopCloser(bytes.NewBufferString(`{
		error:
		}`)),
	}
}

func TestFetchUsageOK(t *testing.T) {
	openaiclient := NewOpenaiClient(setUpMockHttpClient(0, usageStatusOK))
	usage, err := openaiclient.fetchUsage(context.Background(), "mock_access_token")
	assert.NoError(t, err, "expexted nil err")
	assert.Equal(t, 10.0, usage, "expected usage is 10.0")
}

func TestFetchStatusUnauthorized(t *testing.T) {
	openaiclient := NewOpenaiClient(setUpMockHttpClient(0, StatusUnauthorized))
	usage, err := openaiclient.fetchUsage(context.Background(), "mock_access_token")
	assert.ErrorIs(t, err, ErrUnauthorized, "Unauthorized")
	assert.Equal(t, 0.0, usage, "expected usage is 0")

}

func refreshStatusOK(req *http.Request) *http.Response {

	if req.Method == "GET" {
		usage := Usage{RateLimit: RateLimitInfo{
			PrimaryWindow: RateLimitWindow{
				UsedPercent: 10.0,
			},
		}}
		bytesData, _ := json.Marshal(usage)
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(bytesData)),
		}
	}
	refresh := RefreshResp{
		AccessToken:  "123",
		ExpiresIn:    123,
		IdToken:      "123",
		RefreshToken: "123",
		Scope:        "123",
		TokenType:    "123",
	}
	bytesData, _ := json.Marshal(refresh)
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(bytesData)),
	}
}

func TestFetchTokenOK(t *testing.T) {
	openaiclient := NewOpenaiClient(setUpMockHttpClient(0, refreshStatusOK))
	access, refresh, _, err := openaiclient.fetchToken(context.Background(), "mock_access_token")
	assert.NoError(t, err, "expexted nil err")
	assert.Equal(t, "123", access, "expected usage is 123")
	assert.Equal(t, "123", refresh, "expected usage is 123")
	// assert.Equal(t, 123, exp, "expected usage is 123")
}

func TestFetchTokenStatusUnauthorized(t *testing.T) {
	openaiclient := NewOpenaiClient(setUpMockHttpClient(0, StatusUnauthorized))
	access, refresh, _, err := openaiclient.fetchToken(context.Background(), "mock_access_token")
	assert.ErrorIs(t, err, ErrUnauthorized, "Unauthorized")
	assert.Equal(t, "", access, "expected usage is 123")
	assert.Equal(t, "", refresh, "expected usage is 123")

}

func TestNewcodexProxyService(t *testing.T) {
	config := &CodexProxyConfig{
		RelayClient: setUpMockHttpClient(0, mockStatusOK),
		FetchClient: setUpMockHttpClient(0, mockStatusOK),
	}
	codexProxyService, _ := NewcodexProxyService(config)
	assert.NotNil(t, codexProxyService, "codexProxyService should not be nil")
}

func TestServiceAddAccount(t *testing.T) {
	accessToken, _ := utils.GenAccessToken(time.Now().Add(time.Hour).Unix(), time.Now().Unix())
	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			FetchClient: setUpMockHttpClient(0, usageStatusOK),
			RelayClient: setUpMockHttpClient(0, mockStatusOK),
		})
	assert.Nil(t, err)
	err = service.AddAccount("test_acc", "", accessToken, "refresh")
	assert.NoError(t, err, "expexted nil error,but get error:%v")
}

// NewcodexProxyService(setUpMockHttpClient())
func ServiceAddAccountClient(req *http.Request) *http.Response {
	if req.Method == "GET" { //mock get usage
		return usageStatusOK(req)
	} else if req.Host == "" {
		return nil
	}
	return nil
}

func mockStatusOK(req *http.Request) *http.Response {

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

}

func mockStatusUnauthorized(req *http.Request) *http.Response {

	return &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

}

func mockStatusTooManyRequests(req *http.Request) *http.Response {

	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

}

func TestServiceAddBatchAccounts(t *testing.T) {
	accessToken, _ := utils.GenAccessToken(time.Now().Add(time.Hour).Unix(), time.Now().Unix())

	mockJSON := []byte(fmt.Sprintf(`{
        "type": "Codex",
        "accounts": [
            {
                "name": "test_acc_1",
                "access_token": "%s",
                "refresh_token": "fake_refresh_token_1"
            },
            {
                "name": "test_acc_2",
                "access_token": "%s",
                "refresh_token": "fake_refresh_token_2"
            }
        ]
    }`, accessToken, accessToken))

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	err := os.WriteFile(tempFile, mockJSON, 0644)

	require.Nil(t, err)
	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			FetchClient: setUpMockHttpClient(0, usageStatusOK),
			RelayClient: setUpMockHttpClient(0, mockStatusOK),
			MaxRetry:    5,
		})
	assert.Nil(t, err)
	num, err := service.AddBatchAccounts(tempFile)
	assert.Nil(t, err)
	assert.Equal(t, 2, num)
	assert.Equal(t, 2, len(service.pool.accountsMap))
}

func TestServiceSavetoDisk(t *testing.T) {
	accessToken, _ := utils.GenAccessToken(time.Now().Add(time.Hour).Unix(), time.Now().Unix())
	mockJSON := []byte(fmt.Sprintf(`{
        "type": "Codex",
        "accounts": [
            {
                "name": "test_acc_1",
                "access_token": "%s",
                "refresh_token": "fake_refresh_token_1"
            }
        ]
    }`, accessToken))

	mockJSONRefreshed := []byte(`{
        "type": "Codex",
        "accounts": [
            {
                "name": "test_acc_1",
                "access_token": "123",
                "refresh_token": "123"
            }
        ]
    }`)

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")
	_ = os.WriteFile(tempFile, mockJSON, 0644)

	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			AccountsFile: tempFile,
			FetchClient:  setUpMockHttpClient(0, refreshStatusOK),
			RelayClient:  setUpMockHttpClient(0, mockStatusUnauthorized),
			MaxRetry:     5,
		})
	require.Nil(t, err)
	n, err := service.AddBatchAccounts(service.config.AccountsFile)
	require.Equal(t, 1, n)
	assert.Nil(t, err)
	resp, err := service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected not found")
	time.Sleep(time.Millisecond * 200)
	byetsData, err := os.ReadFile(tempFile)
	require.Nil(t, err)
	var c1, c2 AccountJsonFile
	_ = json.Unmarshal(byetsData, &c1)
	_ = json.Unmarshal(mockJSONRefreshed, &c2)
	assert.Equal(t, c1, c2)

}

func TestServiceDoProxyRequestEnabled(t *testing.T) {
	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			FetchClient: setUpMockHttpClient(0, usageStatusOK),
			RelayClient: setUpMockHttpClient(0, mockStatusOK),
			MaxRetry:    5,
		})
	assert.Nil(t, err)
	resp, err := service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrEmptyPool, "expected empty pool")

	_ = service.pool.AddAccount(&Account{
		ID:           "1",
		Status:       Enabled,
		AccessToken:  "123",
		UsagePercent: 10.0,
	})
	resp, err = service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.NotNil(t, resp, "expected  not nil resp")
	assert.NotErrorIs(t, err, ErrEmptyPool, "expected not empty pool")
	assert.Equal(t, resp.StatusCode, 200, "expected status ok")

	account, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	account.UpdateStatus(Disabled)
	resp, err = service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")

	_ = service.pool.AddAccount(&Account{
		ID:           "2",
		Status:       Enabled,
		AccessToken:  "123",
		UsagePercent: 10.0,
	})
	resp, err = service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.NotNil(t, resp, "expected  not nil resp")
	assert.NotErrorIs(t, err, ErrEmptyPool, "expected not empty pool")
	assert.Equal(t, resp.StatusCode, 200, "expected status ok")
}

func TestServiceDoProxyRequestColding(t *testing.T) {

	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			FetchClient: setUpMockHttpClient(0, usageStatusOK),
			RelayClient: setUpMockHttpClient(0, mockStatusTooManyRequests),
			ColdingTime: time.Millisecond * 10,
		})
	assert.Nil(t, err)
	_ = service.pool.AddAccount(&Account{
		ID:           "1",
		Status:       Enabled,
		AccessToken:  "123",
		UsagePercent: 10.0,
	})
	resp, err := service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")
	acc, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	assert.Equal(t, acc.GetStatus(), Colding, "expected status colding")
	time.Sleep(time.Millisecond * 20)
	assert.Equal(t, acc.GetStatus(), Enabled, "expected status colding")

}

func TestServiceDoProxyRequestRefresh(t *testing.T) {

	service, err := NewcodexProxyService(
		&CodexProxyConfig{
			FetchClient: setUpMockHttpClient(time.Millisecond*10, refreshStatusOK),
			RelayClient: setUpMockHttpClient(0, mockStatusUnauthorized),
			ColdingTime: 0,
		})
	assert.Nil(t, err)
	_ = service.pool.AddAccount(&Account{
		ID:           "1",
		Status:       Enabled,
		AccessToken:  "123",
		UsagePercent: 10.0,
	})
	resp, err := service.DoProxyRequest(context.Background(), []byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")
	acc, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	assert.Equal(t, acc.GetStatus(), Refreshing, "expected status refreshing")
	time.Sleep(time.Millisecond * 30)
	assert.Equal(t, acc.GetStatus(), Enabled, "expected status colding")

}
