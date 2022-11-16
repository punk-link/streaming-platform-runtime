package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/punk-link/logger"
	contracts "github.com/punk-link/platform-contracts"
)

func CreateJstStreamIfNotExist(err error, logger logger.Logger, jetStreamContext nats.JetStreamContext) error {
	if err != nil {
		return err
	}

	stream, _ := jetStreamContext.StreamInfo(contracts.PLATFORM_URL_RESPONSE_STREAM_NAME)
	if stream == nil {
		logger.LogInfo("Creating Nats stream %s and subjects %s", contracts.PLATFORM_URL_RESPONSE_STREAM_NAME, contracts.PLATFORM_URL_RESPONSE_STREAM_SUBJECT)
		_, err = jetStreamContext.AddStream(contracts.DefaultReducerConfig)
	}

	return err
}

func GetSubscription(err error, jetStreamContext nats.JetStreamContext, platformName string) (*nats.Subscription, error) {
	if err != nil {
		return nil, err
	}

	subject := contracts.GetRequestStreamSubject(platformName)
	consumerName := contracts.GetRequestConsumerName(platformName)
	return jetStreamContext.PullSubscribe(subject, consumerName)
}
