package processing

import (
	"context"
	"encoding/json"
	"sync"

	natsHelper "github.com/punk-link/streaming-platform-runtime/nats"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	contracts "github.com/punk-link/platform-contracts"
	runtime "github.com/punk-link/streaming-platform-runtime"
)

type QueueProcessingService struct {
	logger             logger.Logger
	natsConnection     *nats.Conn
	urlsInProcess      syncint64.UpDownCounter
	urlsProcessedTotal syncint64.Counter
	urlsReceivedTotal  syncint64.Counter
}

func New(options *runtime.ServiceOptions, natsConnection *nats.Conn) QueueProcessor {
	meter := global.MeterProvider().Meter(options.ServiceName)
	urlsInProcess, _ := meter.SyncInt64().UpDownCounter("release_urls_in_process")
	urlsProcessedTotal, _ := meter.SyncInt64().Counter("urls_processed")
	urlsReceivedTotal, _ := meter.SyncInt64().Counter("urls_received")

	return &QueueProcessingService{
		logger:             options.Logger,
		natsConnection:     natsConnection,
		urlsInProcess:      urlsInProcess,
		urlsProcessedTotal: urlsProcessedTotal,
		urlsReceivedTotal:  urlsReceivedTotal,
	}
}

func (t *QueueProcessingService) Process(ctx context.Context, wg *sync.WaitGroup, platformer contracts.Platformer) {
	defer wg.Done()

	jetStreamContext, err := t.natsConnection.JetStream()
	subscription, err := natsHelper.GetSubscription(err, jetStreamContext, platformer.GetPlatformName())
	if err != nil {
		t.logger.LogError(err, err.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, _ := subscription.Fetch(platformer.GetBatchSize())
			containers := make([]contracts.UpcContainer, len(messages))
			for i, message := range messages {
				var container contracts.UpcContainer
				_ = json.Unmarshal(message.Data, &container)

				containers[i] = container
				message.Ack()
			}

			t.urlsInProcess.Add(ctx, int64(len(containers)))
			t.urlsReceivedTotal.Add(ctx, int64(len(containers)))

			urlResults := platformer.GetReleaseUrlsByUpc(containers)
			err = natsHelper.CreateJstStreamIfNotExist(nil, t.logger, jetStreamContext)
			_ = t.publishUrlResults(err, ctx, jetStreamContext, urlResults)
		}
	}
}

func (t *QueueProcessingService) publishUrlResults(err error, ctx context.Context, jetStreamContext nats.JetStreamContext, urlResults []contracts.UrlResultContainer) error {
	if err != nil {
		return err
	}

	for _, urlResult := range urlResults {
		json, _ := json.Marshal(urlResult)
		jetStreamContext.Publish(contracts.PLATFORM_URL_RESPONSE_STREAM_SUBJECT, json)
	}

	t.urlsInProcess.Add(ctx, -int64(len(urlResults)))
	t.urlsProcessedTotal.Add(ctx, int64(len(urlResults)))

	return err
}
