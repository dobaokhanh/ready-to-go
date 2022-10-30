package ratelimiter_leakybucket

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	ratelimit "go.uber.org/ratelimit"
)

func RateLimit(rate int) gin.HandlerFunc {
	var lmap sync.Map

	return func(ctx *gin.Context) {
		ipAddr := ctx.ClientIP()

		lif, ok := lmap.Load(ipAddr)
		if !ok {
			lif = ratelimit.New(rate)
		}

		lm, ok := lif.(ratelimit.Limiter)
		if !ok {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("something wrong happened"))
			return
		}

		lm.Take()
		lmap.Store(ipAddr, lm)

		ctx.Next()
	}
}
