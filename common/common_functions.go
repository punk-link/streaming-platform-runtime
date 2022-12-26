package common

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	"github.com/punk-link/logger"
	"github.com/punk-link/streaming-platform-runtime/constants"
	vaultClient "github.com/punk-link/vault-client"
)

func GetAppSecrets(envManager envManager.EnvironmentVariableManager, logger logger.Logger, storeName string, secretName string) map[string]any {
	vaultAddress, isExist := envManager.TryGet("PNKL_VAULT_ADDR")
	if !isExist {
		err := errors.New("can't get PNKL_VAULT_ADDR environment variable")
		logger.LogFatal(err, err.Error())
	}

	vaultToken, isExist := envManager.TryGet("PNKL_VAULT_TOKEN")
	if !isExist {
		err := errors.New("an't get PNKL_VAULT_TOKEN environment variable")
		logger.LogFatal(err, err.Error())
	}

	vaultConfig := &vaultClient.VaultClientOptions{
		Endpoint: vaultAddress,
		RoleName: secretName,
	}

	vaultClient := vaultClient.New(vaultConfig, logger)
	return vaultClient.Get(vaultToken, storeName, secretName)
}

func GetConsulClient(logger logger.Logger, appSecrets map[string]any, environmentName string, serviceName string) consulClient.ConsulClient {
	consul, err := consulClient.New(&consulClient.ConsulConfig{
		Address:         appSecrets["consul-address"].(string),
		EnvironmentName: environmentName,
		StorageName:     serviceName,
		Token:           appSecrets["consul-token"].(string),
	})
	if err != nil {
		logger.LogFatal(err, err.Error())
	}

	return consul
}

func GetEnvironmentName(envManager envManager.EnvironmentVariableManager) string {
	name, isExist := envManager.TryGet(constants.GO_ENVIRONMENT)
	if !isExist {
		return constants.DEFAULT_GO_ENVIRONMENT
	}

	return name
}

func GetNatsConnection(logger logger.Logger, consul consulClient.ConsulClient) *nats.Conn {
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
