package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/domains/link"
	"gorm.io/gorm"
)

type LinkHandler struct {
	db *gorm.DB
}

func NewLinkHandler(db *gorm.DB) *LinkHandler {
	return &LinkHandler{
		db,
	}
}

func (l *LinkHandler) Get(c *gin.Context) {
}

func (l *LinkHandler) Create(c *gin.Context) {
	path := c.Request.FormValue("path")
	url := c.Request.FormValue("url")

	if path == "" {
		c.String(http.StatusUnprocessableEntity, "path is required")
		return
	}

	if url == "" {
		c.String(http.StatusUnprocessableEntity, "url is required")
		return
	}

	l.db.Create(&link.Link{
		Path: path,
		URL:  url,
	})

	c.Redirect(http.StatusSeeOther, "/")
}
