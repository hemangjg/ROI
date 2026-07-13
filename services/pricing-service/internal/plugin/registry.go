package plugin

import (
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/domain"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin/anthropic"
	"github.com/ai-finops/ai-finops/services/pricing-service/internal/plugin/openai"
)

// Registry holds all provider pricing plugins.
type Registry struct {
	plugins []domain.ProviderPlugin
}

// NewRegistry wires Phase 1 provider plugins.
func NewRegistry() *Registry {
	return &Registry{
		plugins: []domain.ProviderPlugin{
			openai.NewPlugin(),
			anthropic.NewPlugin(),
		},
	}
}

// All returns registered plugins in stable order.
func (r *Registry) All() []domain.ProviderPlugin {
	return r.plugins
}