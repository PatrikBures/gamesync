package middlewares

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ogen-go/ogen/middleware"
)

func Logging(logger *slog.Logger) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		logger := logger.With(
			"operation", req.OperationName,
		)

		start := time.Now()

		resp, err := next(req)
		if err != nil {
			logger.Error("Fail", "error", err)
			return resp, err
		}

		response_time := time.Since(start)

		fields := []any{"response_time", response_time}

		if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
			fields = append(fields, "status_code", tresp.GetStatusCode())
		} else {
			fields = append(fields, "status_code", "unknown")
			fields = append(fields, "response_type", fmt.Sprintf("%T", resp.Type))
		}

		logger.Info("Success", fields...)
		return resp, nil
	}
}
