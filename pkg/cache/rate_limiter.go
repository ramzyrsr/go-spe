package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Rate limit parameters
const (
	BurstLimit      = 10 // 10 requests in 10 seconds
	BurstWindow     = 10 // 10 seconds
	SustainedLimit  = 5  // 5 requests per minute
	SustainedWindow = 60 // 60 seconds
)

// IsRateLimited checks if the IP address has exceeded the rate limit (burst + sustained)
func IsRateLimited(redisClient *redis.Client, ipAddress string) bool {
	// Burst Key
	burstKey := fmt.Sprintf("burst_limit:%s", ipAddress)
	// Sustained Key
	sustainedKey := fmt.Sprintf("sustained_limit:%s", ipAddress)

	// Check burst rate limiting
	burstAllowed := checkBurstLimit(redisClient, burstKey)
	if !burstAllowed {
		// If burst limit exceeded, apply sustained rate limiting
		return checkSustainedLimit(redisClient, sustainedKey)
	}

	return true
}

// checkBurstLimit checks if the burst rate limit is exceeded (e.g., 10 requests in 10 seconds)
func checkBurstLimit(redisClient *redis.Client, burstKey string) bool {
	// Increment the burst counter
	burstCount, err := redisClient.Incr(context.Background(), burstKey).Result()
	if err != nil {
		return false
	}

	// Set the expiration time for the burst limit (after 10 seconds, reset the counter)
	if burstCount == 1 {
		_, err := redisClient.Expire(context.Background(), burstKey, time.Duration(BurstWindow)*time.Second).Result()
		if err != nil {
			return false
		}
	}

	// If the burst count exceeds the limit, deny the request
	if burstCount > int64(BurstLimit) {
		return false
	}

	return true
}

// checkSustainedLimit checks if the sustained rate limit is exceeded (e.g., 5 requests per minute)
func checkSustainedLimit(redisClient *redis.Client, sustainedKey string) bool {
	// Increment the sustained counter
	sustainedCount, err := redisClient.Incr(context.Background(), sustainedKey).Result()
	if err != nil {
		return false
	}

	// Set the expiration time for the sustained limit (after 60 seconds, reset the counter)
	if sustainedCount == 1 {
		_, err := redisClient.Expire(context.Background(), sustainedKey, time.Duration(SustainedWindow)*time.Second).Result()
		if err != nil {
			return false
		}
	}

	// If the sustained count exceeds the limit, deny the request
	if sustainedCount > int64(SustainedLimit) {
		return false
	}

	return true
}
