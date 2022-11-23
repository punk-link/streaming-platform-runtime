package common

import (
	"fmt"

	"github.com/nats-io/nats.go"
	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	"github.com/punk-link/logger"
	"github.com/punk-link/streaming-platform-runtime/constants"
)

func GetConsulClient(logger logger.Logger, environmentName string, serviceName string) *consulClient.ConsulClient {
	isExist, consulAddress := envManager.TryGetEnvironmentVariable(constants.CONSUL_ADDRESS)
	if !isExist {
		err := fmt.Errorf("can't find value of the '%s' environment variable", constants.CONSUL_ADDRESS)
		logger.LogFatal(err, err.Error())
	}

	isExist, consulToken := envManager.TryGetEnvironmentVariable(constants.CONSUL_TOKEN)
	if !isExist {
		err := fmt.Errorf("can't find value of the '%s' environment variable", constants.CONSUL_TOKEN)
		logger.LogFatal(err, err.Error())
	}

	consul, err := consulClient.New(&consulClient.ConsulConfig{
		Address:         consulAddress,
		EnvironmentName: environmentName,
		StorageName:     serviceName,
		Token:           consulToken,
	})
	if err != nil {
		logger.LogFatal(err, err.Error())
	}

	return consul
}

func GetEnvironmentName() string {
	isExist, name := envManager.TryGetEnvironmentVariable(constants.GO_ENVIRONMENT)
	if !isExist {
		return constants.DEFAULT_GO_ENVIRONMENT
	}

	return name
}

func GetNatsConnection(logger logger.Logger, consul *consulClient.ConsulClient) *nats.Conn {
	natsSettingsValues, err := consul.Get("NatsSettings")
	if err != nil {
		err := fmt.Errorf("can't obtain Nats settings from Consul: '%s'", err.Error())
		logger.LogFatal(err, err.Error())
	}
	natsSettings := natsSettingsValues.(map[string]interface{})

	natsConnection, err := nats.Connect(natsSettings["Endpoint"].(string))
	if err != nil {
		err := fmt.Errorf("can't obtain Nats settings from Consul: '%s'", err.Error())
		logger.LogFatal(err, err.Error())
	}

	return natsConnection
}
