package handlers

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sdisaacson/TriggerLab/internal/api"
	"github.com/sdisaacson/TriggerLab/internal/checkers"
	"github.com/sdisaacson/TriggerLab/internal/config"
	"github.com/sdisaacson/TriggerLab/internal/pubsub"
	"github.com/sdisaacson/TriggerLab/internal/storage"
)

type API struct {
	apiVersion
	apiHealth
	apiSession
	apiSettings
	apiMetrics
	apiWebsocket
}

var _ api.ServerInterface = (*API)(nil) // verify that API implements interface

func NewAPI(
	ctx context.Context,
	cfg config.Config,
	rdb *redis.Client,
	stor storage.Storage,
	pub pubsub.Publisher,
	sub pubsub.Subscriber,
	registry *prometheus.Registry,
	version string,
	websocketMetrics websocketMetrics,
) *API {
	var result = API{}

	result.apiHealth.liveChecker = checkers.NewLiveChecker()
	result.apiHealth.readyChecker = checkers.NewReadyChecker(ctx, rdb)

	result.apiSettings.cfg = cfg

	result.apiVersion.version = version

	result.apiSession.storage = stor
	result.apiSession.pub = pub

	result.apiMetrics.registry = registry

	result.apiWebsocket.ctx = ctx
	result.apiWebsocket.cfg = cfg
	result.apiWebsocket.stor = stor
	result.apiWebsocket.pub = pub
	result.apiWebsocket.sub = sub
	result.apiWebsocket.metrics = websocketMetrics

	return &result
}
