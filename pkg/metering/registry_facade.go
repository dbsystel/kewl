package metering

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ MetricsRegistry = &prometheusRegistryFacade{}

// prometheusRegistryFacade facades a prometheus.Registry providing Summary by alias
type prometheusRegistryFacade struct {
	*prometheus.Registry
	rwMutex sync.RWMutex
	aliases map[string]Summary
}

func (p *prometheusRegistryFacade) WithAliasOrCreate(alias string, promOpts *prometheus.SummaryOpts, labelNames ...string) Summary {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	if existing := p.unsafeWithAlias(alias); existing != nil {
		return existing
	}
	return p.unsafeNewSummary(alias, promOpts, labelNames)
}

func (p *prometheusRegistryFacade) WithAlias(alias string) Summary {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.unsafeWithAlias(alias)
}

func (p *prometheusRegistryFacade) NewSummary(alias string, promOpts *prometheus.SummaryOpts, labelNames ...string) Summary {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	return p.unsafeNewSummary(alias, promOpts, labelNames)
}

func (p *prometheusRegistryFacade) unsafeWithAlias(alias string) Summary {
	if p.aliases == nil {
		return nil
	}
	return p.aliases[alias]
}

func (p *prometheusRegistryFacade) unsafeNewSummary(alias string, promOpts *prometheus.SummaryOpts, labelNames []string) Summary {
	if p.aliases == nil {
		p.aliases = make(map[string]Summary)
	}
	if _, ok := p.aliases[alias]; ok {
		panic("alias already registered: " + alias)
	}
	result := p.createSummary(promOpts, labelNames)
	p.aliases[alias] = result
	return result
}

func (p *prometheusRegistryFacade) createSummary(promOpts *prometheus.SummaryOpts, labelNames []string) Summary {
	factory := promauto.With(p.Registry)
	return &summaryImpl{
		summary: factory.NewSummaryVec(*promOpts, labelNames),
	}
}
