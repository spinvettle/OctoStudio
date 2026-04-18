package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/spinvettle/OctoStudio/internal/utils"
	"github.com/stretchr/testify/assert"
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

func TestNewProxyService(t *testing.T) {
	proxyService := NewProxyService(nil, nil, 0)
	assert.NotNil(t, proxyService, "proxyService should not be nil")
}

func TestServiceAddAccount(t *testing.T) {
	accessToken, _ := utils.GenAccessToken(time.Now().Add(time.Hour).Unix(), time.Now().Unix())
	service := NewProxyService(nil, setUpMockHttpClient(0, usageStatusOK), 0)
	err := service.AddAccount("test_acc", "", accessToken, "refresh")
	assert.NoError(t, err, "expexted nil error,but get error:%v")
}

// NewProxyService(setUpMockHttpClient())
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

func TestServiceDoProxyRequestEnabled(t *testing.T) {
	service := NewProxyService(setUpMockHttpClient(time.Duration(0), mockStatusOK),
		setUpMockHttpClient(time.Duration(0), mockStatusOK), 0)

	resp, err := service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrEmptyPool, "expected empty pool")

	_ = service.pool.AddAccount(&Account{
		ID:          "1",
		Status:      Enabled,
		AccessToken: "123",
	})
	resp, err = service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.NotNil(t, resp, "expected  not nil resp")
	assert.NotErrorIs(t, err, ErrEmptyPool, "expected not empty pool")
	assert.Equal(t, resp.StatusCode, 200, "expected status ok")

	account, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	account.UpdateStatus(Disabled)
	resp, err = service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")

	_ = service.pool.AddAccount(&Account{
		ID:          "2",
		Status:      Enabled,
		AccessToken: "123",
	})
	resp, err = service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.NotNil(t, resp, "expected  not nil resp")
	assert.NotErrorIs(t, err, ErrEmptyPool, "expected not empty pool")
	assert.Equal(t, resp.StatusCode, 200, "expected status ok")
}

func TestServiceDoProxyRequestColding(t *testing.T) {
	service := NewProxyService(setUpMockHttpClient(0, mockStatusTooManyRequests),
		setUpMockHttpClient(0, mockStatusOK), 0)
	_ = service.pool.AddAccount(&Account{
		ID:          "1",
		Status:      Enabled,
		AccessToken: "123",
	})
	resp, err := service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")
	acc, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	assert.Equal(t, acc.GetStatus(), Colding, "expected status colding")
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, acc.GetStatus(), Enabled, "expected status colding")

}

func TestServiceDoProxyRequestRefresh(t *testing.T) {
	service := NewProxyService(setUpMockHttpClient(0, mockStatusUnauthorized),
		setUpMockHttpClient(time.Millisecond*10, refreshStatusOK), 0)
	_ = service.pool.AddAccount(&Account{
		ID:          "1",
		Status:      Enabled,
		AccessToken: "123",
	})
	resp, err := service.DoProxyRequest([]byte{}, make(map[string][]string))
	assert.Nil(t, resp, "expected nil resp")
	assert.ErrorIs(t, err, ErrNotFoundAvailabelAccount, "expected empty pool")
	acc, err := service.pool.GetAccountById("1")
	assert.Nil(t, err)
	assert.Equal(t, acc.GetStatus(), Refreshing, "expected status refreshing")
	time.Sleep(time.Millisecond * 30)
	assert.Equal(t, acc.GetStatus(), Enabled, "expected status colding")

}
