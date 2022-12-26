package streamingplatformruntime

import (
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"github.com/punk-link/streaming-platform-runtime/common"
)

type ServiceOptions struct {
	Consul          consulClient.ConsulClient
	EnvironmentName string
	Logger          logger.Logger
	ServiceName     string
}

func NewServiceOptions(logger logger.Logger, appSecrets map[string]any, environmentName string, serviceName string) *ServiceOptions {
	return &ServiceOptions{
		Consul:          common.GetConsulClient(logger, appSecrets, environmentName, serviceName),
		EnvironmentName: environmentName,
		Logger:          logger,
		ServiceName:     serviceName,
	}
}
