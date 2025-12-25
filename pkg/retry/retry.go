package retry

import (
	"context"
	"flight-aggregator/pkg/config"
	"fmt"
	"log"
	"time"
)

// Params holds configuration for retry logic
type Params struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

// FromConfig creates Params from config.RetryConfig
func FromConfig(cfg config.RetryConfig) Params {
	return Params{
		MaxAttempts:       cfg.MaxAttempts,
		InitialDelay:      cfg.GetInitialDelay(),
		MaxDelay:          cfg.GetMaxDelay(),
		BackoffMultiplier: cfg.Multiplier,
	}
}

// RetryableWithCheckFunc is a function that can be retried and reports if the error is retryable
type RetryableWithCheckFunc func() (error, bool)

// RetryWithCheck executes a function with exponential backoff, checking if errors are retryable
func RetryWithCheck(ctx context.Context, params Params, fn RetryableWithCheckFunc, operationName string) error {
	var lastErr error
	currentDelay := params.InitialDelay

	for attempt := 1; attempt <= params.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: context cancelled: %w", operationName, ctx.Err())
		default:
		}

		// Execute the function
		err, shouldRetry := fn()
		if err == nil {
			// Success
			if attempt > 1 {
				log.Printf("%s: succeeded on attempt %d/%d", operationName, attempt, params.MaxAttempts)
			}
			return nil
		}

		lastErr = err

		// Don't retry if error is not retryable
		if !shouldRetry {
			log.Printf("%s: non-retryable error on attempt %d: %v", operationName, attempt, err)
			return err
		}

		// Don't retry on last attempt
		if attempt == params.MaxAttempts {
			log.Printf("%s: failed after %d attempts: %v", operationName, params.MaxAttempts, err)
			break
		}

		// Log retry attempt
		log.Printf("%s: attempt %d/%d failed: %v, retrying in %v",
			operationName, attempt, params.MaxAttempts, err, currentDelay)

		// Wait before retrying (with context cancellation support)
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: context cancelled during retry wait: %w", operationName, ctx.Err())
		case <-time.After(currentDelay):
		}

		// Calculate next delay with exponential backoff
		currentDelay = time.Duration(float64(currentDelay) * params.BackoffMultiplier)
		if currentDelay > params.MaxDelay {
			currentDelay = params.MaxDelay
		}
	}

	return fmt.Errorf("%s: failed after %d attempts: %w", operationName, params.MaxAttempts, lastErr)
}
