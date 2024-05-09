package middleware

import (
	"github.com/gin-gonic/gin"
	"tailscale.com/client/tailscale"
	"tailscale.com/tailcfg"
)

const (
	UserProfileKey = "user_profile"
)

func Tailscale(lc *tailscale.LocalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		who, err := lc.WhoIs(c, c.Request.RemoteAddr)
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}

		c.Set(UserProfileKey, who.UserProfile)
		c.Next()
	}
}

func GetUserProfile(c *gin.Context) *tailcfg.UserProfile {
	return c.MustGet(UserProfileKey).(*tailcfg.UserProfile)
}
