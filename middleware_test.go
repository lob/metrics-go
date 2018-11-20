package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/lob/assets-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("sends request duration through Datadog client", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		e := echo.New()
		e.Use(Middleware(metrics))

		e.GET("/", func(c echo.Context) error {
			time.Sleep(time.Duration(testDuration) * time.Millisecond)
			return nil
		})

		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.True(t, testDuration <= mc.duration, "incorrect duration")
	})
}
