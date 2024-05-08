package app

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/assets"
	"github.com/tgiv014/to/handlers"
)

func (a *App) Router() *gin.Engine {
	router := gin.Default()
	router.StaticFS("/static", static.EmbedFolder(assets.Assets, "dist"))
	router.GET("/", handlers.IndexHandler{}.Index)

	links := router.Group("links")
	{
		linkHandler := handlers.NewLinkHandler(a.db)
		links.GET("/", linkHandler.Get)
		links.POST("/", linkHandler.Create)
	}

	return router
}
