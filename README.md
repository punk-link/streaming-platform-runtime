# streaming-platform-runtime
Generic functions, and telemetry and runtime for streaming platform services


## Example

```go
package main

import (
	"main/services"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"

	runtime "github.com/punk-link/streaming-platform-runtime"
	common "github.com/punk-link/streaming-platform-runtime/common"
	"github.com/punk-link/streaming-platform-runtime/startup"
)

func main() {
	logger := logger.New()
	environmentName := common.GetEnvironmentName()
	logger.LogInfo("%s is running as '%s'", "my-service", environmentName)

	serviceOptions := runtime.NewServiceOptions(logger, environmentName, SERVICE_NAME)

	myService := services.NewMyService(logger, httpclient.DefaultConfig(logger))
	go startup.ProcessUrls(serviceOptions, myService)

	startup.RunServer(serviceOptions)
}
```