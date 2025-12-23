package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Cache     CacheConfig     `yaml:"cache"`
	Provider  ProviderConfig  `yaml:"provider"`
	Logging   LoggingConfig   `yaml:"logging"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Scoring   ScoringConfig   `yaml:"scoring"`
	Retry     RetryConfig     `yaml:"retry"`
	MockData  MockDataConfig  `yaml:"mock_data"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	Timeout      string `yaml:"timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

type CacheConfig struct {
	TTL string `yaml:"ttl"`
}

type ProviderConfig struct {
	Timeout   string                    `yaml:"timeout"`
	Providers map[string]ProviderDetail `yaml:"providers"`
}

type ProviderDetail struct {
	Name         string `yaml:"name"`
	Enabled      bool   `yaml:"enabled"`
	ResponseTime string `yaml:"response_time"`
	//ResponseTimeEndRange  int `yaml:"response_time_end_range"` //Real world simulation
	FailureRate float64 `yaml:"failure_rate"`
	DataPath    string  `yaml:"data_path"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type RateLimitConfig struct {
	Requests int    `yaml:"requests"`
	Window   string `yaml:"window"`
}

type ScoringConfig struct {
	Weights ScoringWeights `yaml:"weights"`
}

type ScoringWeights struct {
	Price         float64 `yaml:"price"`
	Duration      float64 `yaml:"duration"`
	Stops         float64 `yaml:"stops"`
	DepartureTime float64 `yaml:"departure_time"`
}

type RetryConfig struct {
	MaxAttempts  int     `yaml:"max_attempts"`
	InitialDelay string  `yaml:"initial_delay"`
	MaxDelay     string  `yaml:"max_delay"`
	Multiplier   float64 `yaml:"multiplier"`
}

type MockDataConfig struct {
	Path string `yaml:"path"`
}

// Load reads configuration from .env.yaml file
func Load() (*Config, error) {
	data, err := os.ReadFile(".env.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read .env.yaml: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse .env.yaml: %w", err)
	}

	return &config, nil
}

// Helper methods to parse durations
func (s *ServerConfig) GetTimeout() time.Duration {
	d, _ := time.ParseDuration(s.Timeout)
	return d
}

func (s *ServerConfig) GetReadTimeout() time.Duration {
	d, _ := time.ParseDuration(s.ReadTimeout)
	return d
}

func (s *ServerConfig) GetWriteTimeout() time.Duration {
	d, _ := time.ParseDuration(s.WriteTimeout)
	return d
}

func (s *ServerConfig) GetIdleTimeout() time.Duration {
	d, _ := time.ParseDuration(s.IdleTimeout)
	return d
}

func (c *CacheConfig) GetTTL() time.Duration {
	d, _ := time.ParseDuration(c.TTL)
	return d
}

func (p *ProviderConfig) GetTimeout() time.Duration {
	d, _ := time.ParseDuration(p.Timeout)
	return d
}

func (r *RateLimitConfig) GetWindow() time.Duration {
	d, _ := time.ParseDuration(r.Window)
	return d
}

func (pd *ProviderDetail) GetResponseTime() time.Duration {
	d, _ := time.ParseDuration(pd.ResponseTime)
	return d
}

func (r *RetryConfig) GetInitialDelay() time.Duration {
	d, _ := time.ParseDuration(r.InitialDelay)
	return d
}

func (r *RetryConfig) GetMaxDelay() time.Duration {
	d, _ := time.ParseDuration(r.MaxDelay)
	return d
}

// GetProviderConfig returns configuration for a specific provider by key
func (p *ProviderConfig) GetProviderConfig(key string) (*ProviderDetail, bool) {
	detail, exists := p.Providers[key]
	return &detail, exists
}
