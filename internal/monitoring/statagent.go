package monitoring

import (
	"final-project/internal/config"
)

type Agent struct {
	configuration config.AgentConfig
}

func NewAgent(cfg config.AgentConfig) *Agent {
	return &Agent{cfg}
}
