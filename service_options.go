package streamingplatformruntime

import (
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"github.com/punk-link/streaming-platform-runtime/common"
)

type ServiceOptions struct {
	Consul          *consulClient.ConsulClient
	EnvironmentName string
	Logger          logger.Logger
	ServiceName     string
}

func NewServiceOptions(logger logger.Logger, environmentName string, serviceName string) *ServiceOptions {
	return &ServiceOptions{
		Consul:          common.GetConsulClient(logger, environmentName, serviceName),
		EnvironmentName: environmentName,
		Logger:          logger,
		ServiceName:     serviceName,
	}
}
