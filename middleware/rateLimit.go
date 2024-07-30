package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-rest-api/internal/config"
	"go-rest-api/pkg"
	"go-rest-api/pkg/tools"
	"net/http"
	"os"
	"strconv"
	"time"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := tools.GetClientIP(c.Request)
		limitKey := fmt.Sprintf("rate_limit:%s", ip)
		limit, err := config.Cache.Get(limitKey)
		if err != nil && !errors.Is(err, redis.Nil) {
			pkg.ApiResponse(c, http.StatusInternalServerError, "internal_server_error", "")
			c.Abort()
			return
		} else if limit == "" {
			err = config.Cache.Set(limitKey, 1, time.Minute)
			if err != nil {
				pkg.ApiResponse(c, http.StatusInternalServerError, "internal_server_error", "")
				c.Abort()
				return
			}
		} else {
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				pkg.ApiResponse(c, http.StatusInternalServerError, "internal_server_error", "")
				c.Abort()
				return
			}

			maxApiRequestInt, err := strconv.Atoi(os.Getenv("MAX_API_REQUEST"))
			if err != nil {
				pkg.ApiResponse(c, http.StatusInternalServerError, "internal_server_error", "")
				c.Abort()
				return
			}

			if limitInt >= maxApiRequestInt {
				pkg.ApiResponse(c, http.StatusTooManyRequests, "too_many_request", "")
				c.Abort()
				return
			} else {
				err = config.Cache.Set(limitKey, limitInt+1, time.Minute)
				if err != nil {
					pkg.ApiResponse(c, http.StatusInternalServerError, "internal_server_error", "")
					c.Abort()
					return
				}
			}
		}

		//rl := ratelimit.New(100) // per second
		//now := time.Now()
		//rl.Take()
		//if time.Since(now) <= time.Minute {
		//	pkg.ApiResponse(c, http.StatusTooManyRequests, "too_many_request", "")
		//	fmt.Println("too many requests, abort the request")
		//	c.Abort()
		//	return
		//}

		c.Next()
	}
}
