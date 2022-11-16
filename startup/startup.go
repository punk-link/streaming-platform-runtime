package startup

import (
	"github.com/gin-gonic/gin"
	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Configure(logger logger.Logger, consul *consulClient.ConsulClient, options *StartupOptions) *gin.Engine {
	gin.SetMode(options.GinMode)
	app := gin.Default()

	app.Use(otelgin.Middleware(options.ServiceName))

	app.LoadHTMLGlob("./var/www/templates/**/*.go.tmpl")
	app.Static("/assets", "./var/www/assets")

	//initSentry(app, logger, consul, options.EnvironmentName)
	configureOpenTelemetry(logger, consul, options)
	//setupRouts(app, diContainer)

	return app
}
