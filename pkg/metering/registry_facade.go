package metering

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ MetricsRegistry = &promRegistryFacade{}

// promRegistryFacade facades a prometheus.Registry providing Summary by alias
type promRegistryFacade struct {
	*prometheus.Registry
	rwMutex sync.RWMutex
	aliases map[string]Summary
}

func (p *promRegistryFacade) PromRegistry() *prometheus.Registry {
	return p.Registry
}

func (p *promRegistryFacade) WithAliasOrCreate(alias string, opts *prometheus.SummaryOpts, lbls ...string) Summary {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	if existing := p.unsafeWithAlias(alias); existing != nil {
		return existing
	}
	return p.unsafeNewSummary(alias, opts, lbls)
}

func (p *promRegistryFacade) WithAlias(alias string) Summary {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.unsafeWithAlias(alias)
}

func (p *promRegistryFacade) NewSummary(alias string, promOpts *prometheus.SummaryOpts, labelNames ...string) Summary {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	return p.unsafeNewSummary(alias, promOpts, labelNames)
}

func (p *promRegistryFacade) unsafeWithAlias(alias string) Summary {
	if p.aliases == nil {
		return nil
	}
	return p.aliases[alias]
}

func (p *promRegistryFacade) unsafeNewSummary(alias string, promOpts *prometheus.SummaryOpts, lbls []string) Summary {
	if p.aliases == nil {
		p.aliases = make(map[string]Summary)
	}
	if _, ok := p.aliases[alias]; ok {
		panic("alias already registered: " + alias)
	}
	result := p.createSummary(promOpts, lbls)
	p.aliases[alias] = result
	return result
}

func (p *promRegistryFacade) createSummary(promOpts *prometheus.SummaryOpts, labelNames []string) Summary {
	factory := promauto.With(p.Registry)
	return &summaryImpl{
		summary: factory.NewSummaryVec(*promOpts, labelNames),
	}
}
