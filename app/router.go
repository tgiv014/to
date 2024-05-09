package app

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/assets"
	"github.com/tgiv014/to/handlers"
	"github.com/tgiv014/to/middleware"
)

// Router returns a new gin router
func (a *App) Router() *gin.Engine {
	linkHandler := handlers.NewLinkHandler(a.links)

	router := gin.Default()
	router.Use(middleware.Tailscale(a.lc))
	router.StaticFS("/static", static.EmbedFolder(assets.Assets, "dist"))
	router.GET("/", linkHandler.Index)
	router.GET("/:path", linkHandler.Follow)

	links := router.Group("links")
	{
		// links.GET("/", linkHandler.Get)
		links.POST("/", linkHandler.Create)
		links.DELETE("/:path", linkHandler.Delete)
	}

	return router
}
