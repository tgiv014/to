package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/components"
)

type IndexHandler struct {
}

func (i IndexHandler) Index(c *gin.Context) {
	components.Index().Render(c, c.Writer)
}
