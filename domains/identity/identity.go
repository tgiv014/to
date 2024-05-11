package identity

import "github.com/gin-gonic/gin"

type User struct {
	Name          string
	Login         string
	ProfilePicURL string
}

type Identifier interface {
	Identify(c *gin.Context) (User, error)
}
