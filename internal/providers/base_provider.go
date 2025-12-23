package providers

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// ProviderConfig contains configuration for a provider
type ProviderConfig struct {
	Name         string
	ResponseTime time.Duration
	FailureRate  float64
	DataPath     string
}

// BaseProvider contains common functionality for all providers
type BaseProvider struct {
	name          string
	responseDelay time.Duration
	failureRate   float64 // 0.0 to 1.0 (0% to 100%)
	mockDataPath  string
}

// NewBaseProviderFromConfig creates a new BaseProvider from config
func NewBaseProviderFromConfig(cfg ProviderConfig) BaseProvider {
	return BaseProvider{
		name:          cfg.Name,
		responseDelay: cfg.ResponseTime,
		failureRate:   cfg.FailureRate,
		mockDataPath:  cfg.DataPath,
	}
}

// Name returns the provider name
func (b *BaseProvider) Name() string {
	return b.name
}

// HealthCheck returns true if the provider is healthy
func (b *BaseProvider) HealthCheck() bool {
	// Simple health check - in real implementation, this would check actual API availability
	return rand.Float64() > b.failureRate
}

// SimulateDelay simulates network delay
func (b *BaseProvider) SimulateDelay(ctx context.Context) error {
	if b.responseDelay > 0 {
		select {
		case <-time.After(b.responseDelay):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// SimulateFailure randomly fails based on failure rate
func (b *BaseProvider) SimulateFailure() error {
	if b.failureRate > 0 && rand.Float64() < b.failureRate {
		return fmt.Errorf("%s: %w (simulated failure for testing)", b.name, ErrProviderUnavailable)
	}
	return nil
}
