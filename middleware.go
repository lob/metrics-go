package metrics

import (
	"fmt"

	"github.com/labstack/echo"
)

// Middleware returns an Echo middleware function that begins a timer before a
// request is handled and ends afterwards.
func Middleware(m Metrics) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			methodTag := fmt.Sprintf("method:%s", c.Request().Method)

			t := m.NewTimer("http.request", methodTag)

			if err := next(c); err != nil {
				c.Error(err)
			}

			statusCodeTag := fmt.Sprintf("status_code:%d", c.Response().Status)
			pathTag := fmt.Sprintf("path:%s", c.Path())

			t.End(statusCodeTag, pathTag)

			return nil
		}
	}
}
