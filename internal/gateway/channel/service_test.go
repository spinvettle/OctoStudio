package channel

// import (
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func setUpDB() *gorm.DB {
// 	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
// 	db.AutoMigrate(
// 		&Channel{}, &ChannelKey{},
// 	)
// 	return db
// }

// func setUpMockHttpClient() *httptest.Server {
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var payload RefreshRequestPayload
// 		body, _ := io.ReadAll(r.Body)
// 		_ = json.Unmarshal(body, &payload)
// 		w.Header().Set("Content-Type", "application/json")
// 		switch payload.RefreshToken {
// 		case "200":
// 			resp := RefreshResp{
// 				AccessToken:  "newAccessToken",
// 				ExpiresIn:    int(time.Now().Add(time.Hour).Unix()),
// 				RefreshToken: "newRefreshToken",
// 			}
// 			byteData, _ := json.Marshal(resp)
// 			w.WriteHeader(http.StatusOK)
// 			w.Write(byteData)

// 		case "401":

// 			byteData, _ := json.Marshal("{error:Unauthorized key}")
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write(byteData)
// 		default:
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("error:mock error"))

// 		}
// 	}))
// 	return server
// }

// func TestChannelService_RefreshChannelKey(t *testing.T) {
// 	server := setUpMockHttpClient()
// 	db := setUpDB()
// 	defer server.Close()
// 	channelSvc := NewChannelService(db,
// 		server.Client(),
// 		&ChannelServiceConfig{FetchTimeout: time.Second * 10})

// 	t.Run("success updates token and refresh token", func(t *testing.T) {
// 		// ...
// 		ck := &ChannelKey{
// 			ID:        1,
// 			ChannelID: 1,
// 			Name:      "testKey",
// 			ApiKey:    "successKeyCase",
// 			Metadata: ChannelKeyMetaData{
// 				RefreshBaseURL:        server.URL,
// 				RefreshRequestPayload: RefreshRequestPayload{RefreshToken: "200"},
// 			},
// 			Status: Enable,
// 		}
// 		n, err := channelSvc.repo.AddChannelKey(ck)
// 		require.Nil(t, err)
// 		require.Equal(t, n, 1)
// 		_, err = channelSvc.RefreshKeyOnce(ck)
// 		assert.Nil(t, err)
// 		ch, err := channelSvc.repo.GetChannelKeyByID(ck.ID)
// 		require.Nil(t, err)
// 		token := ch.Metadata.RefreshRequestPayload.RefreshToken
// 		assert.Equal(t, token,
// 			"newRefreshToken",
// 			"expected \"newRefreshToken\" but get %s", token)

// 	})

// 	t.Run("returns unauthorized error on 401", func(t *testing.T) {
// 		ck := &ChannelKey{
// 			ID:        2,
// 			ChannelID: 2,
// 			Name:      "testKey",
// 			ApiKey:    "successKeyCase",
// 			Metadata: ChannelKeyMetaData{
// 				RefreshBaseURL:        server.URL,
// 				RefreshRequestPayload: RefreshRequestPayload{RefreshToken: "401"},
// 			},
// 			Status: Enable,
// 		}
// 		n, err := channelSvc.repo.AddChannelKey(ck)
// 		require.Nil(t, err)
// 		require.Equal(t, n, 1)
// 		_, err = channelSvc.RefreshKeyOnce(ck)
// 		assert.NotNil(t, err)
// 		require.Equal(t, ErrRefreshTokenUnauthorized, err,
// 			"expected ErrRefreshTokenUnauthorized but get %v", err)
// 		ch, err := channelSvc.repo.GetChannelKeyByID(ck.ID)
// 		require.Nil(t, err)
// 		require.Equal(t, Disable, ch.Status, "expected key suatus is Disable but get %s", ch.Status)

// 	})

// 	t.Run("concurrent refreshing", func(t *testing.T) {
// 		t.Parallel()
// 		var wg sync.WaitGroup
// 		ck := &ChannelKey{
// 			ID:        3,
// 			ChannelID: 3,
// 			Name:      "testKey",
// 			ApiKey:    "successKeyCase",
// 			Metadata: ChannelKeyMetaData{
// 				RefreshBaseURL:        server.URL,
// 				RefreshRequestPayload: RefreshRequestPayload{RefreshToken: "401"},
// 			},
// 			Status: Enable,
// 		}
// 		n, err := channelSvc.repo.AddChannelKey(ck)
// 		require.Nil(t, err)
// 		require.Equal(t, n, 1)
// 		var shared bool
// 		for i := 0; i < 5; i++ {
// 			wg.Add(1)
// 			go func() {
// 				shared, err = channelSvc.RefreshKeyOnce(ck)
// 				wg.Done()
// 			}()
// 		}
// 		assert.Equal(t, shared, true)
// 		assert.NotNil(t, err)
// 		require.Nil(t, err)
// 		ch, err := channelSvc.repo.GetChannelKeyByID(ck.ID)
// 		require.Equal(t, Disable, ch.Status, "expected key suatus is Disable but get %s", ch.Status)

// 		// ...
// 	})
// }
