package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	contracts "github.com/punk-link/platform-contracts"
	runtime "github.com/punk-link/streaming-platform-runtime"
	"github.com/punk-link/streaming-platform-runtime/common"
	"github.com/punk-link/streaming-platform-runtime/processing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func ProcessUrls(options *runtime.ServiceOptions, service contracts.Platformer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	natsConnection := common.GetNatsConnection(options.Logger, options.Consul)
	queueProcessingService := processing.New(options, natsConnection)
	go queueProcessingService.Process(ctx, &wg, service)

	wg.Wait()
	options.Logger.LogInfo("Exiting...")
}

func RunServer(options *runtime.ServiceOptions) {
	options.Logger.LogInfo("Starting Gin server...")

	hostSettingsValues, err := options.Consul.Get("HostSettings")
	if err != nil {
		options.Logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	hostSettings := hostSettingsValues.(map[string]any)

	gin.SetMode(hostSettings["Mode"].(string))
	app := gin.Default()

	app.Use(otelgin.Middleware(options.ServiceName))

	initSentry(app, options)
	configureOpenTelemetry(options)
	setupRouts(app)

	err = app.Run(fmt.Sprintf(":%s", hostSettings["Port"]))
	if err != nil {
		options.Logger.LogFatal(err, fmt.Sprintf("Can't run Gin server: %s", err.Error()))
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostSettings["Address"], hostSettings["Port"]),
		Handler: app,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			options.Logger.LogError(err, "Server listen error: %s\n", err.Error())
		}
	}()
}
