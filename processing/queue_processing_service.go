package processing

import (
	"context"
	"encoding/json"
	"sync"

	natsHelper "github.com/punk-link/streaming-platform-runtime/nats"

	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	contracts "github.com/punk-link/platform-contracts"
)

type QueueProcessingService struct {
	logger         logger.Logger
	natsConnection *nats.Conn
}

func New(logger logger.Logger, natsConnection *nats.Conn) *QueueProcessingService {
	return &QueueProcessingService{
		logger:         logger,
		natsConnection: natsConnection,
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
				message.Ack()

				var container contracts.UpcContainer
				_ = json.Unmarshal(message.Data, &container)

				containers[i] = container
			}

			urlResults := platformer.GetReleaseUrlsByUpc(containers)
			err = natsHelper.CreateJstStreamIfNotExist(nil, t.logger, jetStreamContext)
			_ = t.publishUrlResults(err, jetStreamContext, urlResults)
		}
	}
}

func (t *QueueProcessingService) publishUrlResults(err error, jetStreamContext nats.JetStreamContext, urlResults []contracts.UrlResultContainer) error {
	if err != nil {
		return err
	}

	for _, urlResult := range urlResults {
		json, _ := json.Marshal(urlResult)
		jetStreamContext.Publish(contracts.PLATFORM_URL_RESPONSE_STREAM_SUBJECT, json)
	}

	return err
}
