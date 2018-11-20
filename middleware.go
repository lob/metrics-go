package metrics

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/lob/assets-proxy/pkg/errutil"
)

// Middleware returns an Echo middleware function that begins a timer before a
// request is handled and ends afterwards.
func Middleware(m Metrics) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			methodTag := fmt.Sprintf("method:%s", c.Request().Method)

			t := m.NewTimer("http.request", methodTag)

			if err := next(c); err != nil {
				if errutil.IsIgnorableErr(err) {
					// Metrics is the first middleware that's registered with the Echo framework. We
					// cannot return the error here as it will be bubbled up to Echo's default error
					// handler which we are trying to prevent from calling in this situation.
					return nil
				}
				c.Error(err)
			}

			statusCodeTag := fmt.Sprintf("status_code:%d", c.Response().Status)
			pathTag := fmt.Sprintf("path:%s", c.Path())

			t.End(statusCodeTag, pathTag)

			return nil
		}
	}
}
