package api

import (
	"github.com/StanislavIvanovQA/simple_bank/config"
	db "github.com/StanislavIvanovQA/simple_bank/db/sqlc"
	"github.com/StanislavIvanovQA/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	cfg := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, cfg)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
