package metrics

import (
	"fmt"

	"github.com/labstack/echo"
)

// Metrics defines the interface for metrics reporters
type Metrics interface {
	Count(name string, count int64, tags ...string)
	Histogram(name string, value float64, tags ...string)
	NewTimer(name string, tags ...string) Timer
}

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
