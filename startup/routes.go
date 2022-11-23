package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/punk-link/streaming-platform-runtime/controllers"
)

func setupRouts(app *gin.Engine) {
	app.GET("/metrics", controllers.GetMetrics)
	app.GET("/health", controllers.CheckHealth)
}
