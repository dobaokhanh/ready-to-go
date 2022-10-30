package ratelimiter

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// RateLimiterConfig holds the basic config we need to create a middleware http.Handler object that
// performs rate limiting before offloading the request to an actual handler.
type RateLimiterConfig struct {
	Strategy    Strategy
	Expiration  time.Duration
	MaxRequests uint64
}

func UseRateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		result, err := config.Strategy.Run(ctx, &Request{
			Key:      key,
			Limit:    config.MaxRequests,
			Duration: config.Expiration,
		})

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "failed to run rate limiting for request: %v", err))
			ctx.Next()
		}

		// when the state is Deny, just return a 429 response to the client and stop the request handling flow
		if result.State == Deny {
			ctx.AbortWithError(http.StatusTooManyRequests, errors.New("Too many requests!"))
			ctx.Next()
		}
	}
}
