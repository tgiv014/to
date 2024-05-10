package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/components"
	"github.com/tgiv014/to/domains/link"
	"github.com/tgiv014/to/middleware"
)

type LinkHandler struct {
	links *link.Service
}

func NewLinkHandler(links *link.Service) *LinkHandler {
	return &LinkHandler{
		links: links,
	}
}

func (l *LinkHandler) Index(c *gin.Context) {
	user := middleware.GetUserProfile(c)
	log.Info("user", user)

	links, err := l.links.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	components.Index(links).Render(c, c.Writer)
}

func (l *LinkHandler) Follow(c *gin.Context) {
	path := c.Param("path")

	link, err := l.links.GetByPath(path)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, link.URL)
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

	err := l.links.Create(&link.Link{
		Path: path,
		URL:  url,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (l *LinkHandler) Delete(c *gin.Context) {
	path := c.Param("path")

	link, err := l.links.GetByPath(path)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = l.links.Delete(link)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
