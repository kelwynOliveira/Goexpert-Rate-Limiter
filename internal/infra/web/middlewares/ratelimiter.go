package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/logger"
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/ratelimiter"
	limiter "github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/ratelimiter/strategies"
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/responsehandler"
)

type RateLimiterMiddlewareInterface interface {
	Handle(next http.Handler) http.Handler
}

type RateLimiterMiddleware struct {
	Logger          zerolog.Logger
	ResponseHandler responsehandler.WebResponseHandlerInterface
	Limiter         ratelimiter.RateLimiterInterface
}

func NewRateLimiterMiddleware(
	logger logger.LoggerInterface,
	responseHandler responsehandler.WebResponseHandlerInterface,
	limiter ratelimiter.RateLimiterInterface,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		Logger:          logger.GetLogger(),
		ResponseHandler: responseHandler,
		Limiter:         limiter,
	}
}

func (rlm *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := rlm.Limiter.Check(r.Context(), r)
		if err != nil {
			rlm.Logger.Error().Err(err).Msg("Rrror while checking rate limit")

			rlm.ResponseHandler.RespondWithError(
				w,
				http.StatusInternalServerError,
				errors.Join(errors.New("error while checking rate limit"), err),
			)

			return
		}

		writeHeaders(w, result)

		rlm.Logger.Debug().Msgf("Rate limit result: %+v", result)

		if result.Result == limiter.Deny {
			rlm.ResponseHandler.RespondWithError(w, http.StatusTooManyRequests, errors.New("rate limit exceeded"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeHeaders(w http.ResponseWriter, res *limiter.RateLimiterResult) {
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(res.Limit, 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(res.Remaining, 10))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(res.ExpiresAt.Unix(), 10))
}
