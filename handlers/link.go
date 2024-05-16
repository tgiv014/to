package handlers

import (
	"fmt"
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
	user := middleware.GetUser(c)
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

	err = l.links.Follow(link)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, link.URL)
}

func (l *LinkHandler) Create(c *gin.Context) {
	path := c.Request.FormValue("path")
	url := c.Request.FormValue("url")

	newLink, err := link.New(path, url)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = l.links.Create(newLink)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (l *LinkHandler) Preview(c *gin.Context) {
	path := c.Query("path")
	url := c.Query("url")

	newLink, err := link.New(path, url)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(c, c.Writer)
		return
	}

	components.PreviewMessage(fmt.Sprintf("%s will redirect to %s", newLink.Path, newLink.URL)).Render(c, c.Writer)
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

	c.Status(http.StatusOK)
}
