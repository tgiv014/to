package handlers

import (
	"bytes"
	"fmt"
	"io"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tgiv014/to/components"
)

type HotReloadHandler struct {
	uuid uuid.UUID
}

func NewHotReloadHandler() *HotReloadHandler {
	return &HotReloadHandler{uuid.New()}
}

func (h *HotReloadHandler) ShouldReload(c *gin.Context) {
	clientUUID := c.DefaultQuery("uuid", "")
	messageChannel := make(chan string)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	if clientUUID == "" {
		log.Info("Setting client UUID")
		go func() {
			buf := new(bytes.Buffer)
			components.Reloader(fmt.Sprintf("/should-reload?uuid=%s", h.uuid)).Render(c, buf)
			messageChannel <- buf.String()
		}()
	} else if clientUUID != h.uuid.String() {
		go func() {
			buf := new(bytes.Buffer)
			components.Refresher().Render(c, buf)
			messageChannel <- buf.String()
		}()
	}

	// c.Header("hx-redirect", "/")
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-messageChannel; ok {
			c.SSEvent("livereload", msg)
			return true
		}
		return false
	})
}
