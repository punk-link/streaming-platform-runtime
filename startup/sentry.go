package startup

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	runtime "github.com/punk-link/streaming-platform-runtime"
)

func initSentry(app *gin.Engine, options *runtime.ServiceOptions) {
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	dsn, _ := options.Consul.Get("SentryDsn")
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              dsn.(string),
		Environment:      options.EnvironmentName,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		options.Logger.LogError(err, "Sentry initialization failed: %v", err.Error())
	}
}
