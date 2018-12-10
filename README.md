# metrics-go

This package provides an interface for reporting metrics to statsd,
along with Echo middleware for use with HTTP APIs.

## Usage

Include the metrics-go package in your project:

```go
import (
    "github.com/lob/metrics-go"
)
```

Create a new Metrics client:

```go

cfg := metrics.Config{ ... }
m := metrics.New(cfg)
```

All metrics reported from the instance will have the environment, container, and release added as tags.

You may also pass additional tags to all `metrics` methods as variadic string arguments.

You may wish to embed the `metrics.Config` struct in your applications configuration
struct for convenience.

### Count

```go
m.Count("event-counter", 1)
```

### Histogram

```go
m.Histogram("queue-depth", 10)
```

### Timers

```go
t := metrics.NewTimer("api-call")

err := apiCall()
if err != nil {
    // End also takes tags that are added to the timer metric
    t.End("state:failed")
} else {
    t.End("state:success")
}
```

### Middleware

This package includes middleware suitable for use with Echo servers.

```go

m := metrics.New(cfg)
e := echo.New()
e.Use(metrics.Middleware(m))

```

## Development

```
# Install necessary dependencies for development
make setup

# Ensure dependencies are safe, complete, and reproducible
make deps

# Run tests and generate test coverage report
make test

# Run linter
make lint

# Remove temporary files generated by the build command and the build directory
make clean
```
