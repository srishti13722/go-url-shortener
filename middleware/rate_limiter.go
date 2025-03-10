package middleware

import (
	"context"
	"fmt"
	"go-url-shortener/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ulule/limiter/v3"
	limiter_redis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimiter(maxRequest int, duration string) fiber.Handler {
	rateString := fmt.Sprintf("%d-%s", maxRequest, duration)

	// parse the rate limit duration
	rate, err := limiter.NewRateFromFormatted(rateString)
	if err != nil {
		log.Fatalf("Invalid rate limit format: %v", err)
	}

	//create a redis based rate limiter store
	store, err := limiter_redis.NewStoreWithOptions(database.RedisClient, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 1,
	})
	if err != nil {
		fmt.Println("Failed to create rate limiter")
	}

	// create a rate limiter instance
	rateLimiter := limiter.New(store, rate)

	return func(c *fiber.Ctx) error {
		//user ip

		ip := c.IP()

		//check rate limit for the user

		limitCtx, err := rateLimiter.Get(context.Background(), ip)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiter error"})
		}

		if limitCtx.Reached {
			resetTime := time.Until(time.Unix(limitCtx.Reset, 0))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many requests. Try again later",
				"retry_after": resetTime.Seconds(),
			})
		}

		// continue processing the request
		return c.Next()
	}
}
