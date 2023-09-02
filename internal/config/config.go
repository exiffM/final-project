package config

import "time"

type AgentConfig struct {
	Timeout     time.Duration // Time of outputing system data
	AvgTime     time.Duration // Time of accumulating and averaging system data
	StorageTime time.Duration // Time of storaging system data
}

func NewConfig(t, avg, strg time.Duration) *AgentConfig {
	return &AgentConfig{t, avg, strg}
}
