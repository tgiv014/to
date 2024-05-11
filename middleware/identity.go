package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/domains/identity"
)

const (
	UserKey = "user"
)

func Identity(provider identity.Identifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := provider.Identify(c)
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}

		c.Set(UserKey, user)
		c.Next()
	}
}

func GetUser(c *gin.Context) identity.User {
	return c.MustGet(UserKey).(identity.User)
}
