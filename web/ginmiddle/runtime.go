package ginmiddle

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ctxKey = "ctx"
)

// SetStdContext store the standard library context.Context in web.Context.
func SetStdContext(c *gin.Context, ctx context.Context) {
	c.Set(ctxKey, ctx)
}

// GetStdContext get standard library context.Context from web.Context.
func GetStdContext(c *gin.Context) context.Context {
	v, ok := c.Get(ctxKey)
	if !ok {
		panic("no context set, you should use ginmiddle.Trace with *web.Engine.")
	}
	return v.(context.Context)
}

// IsWekSocket check whether the request want to upgrade to websocket.
func IsWekSocket(c *gin.Context) bool {
	connection := c.GetHeader("Connection")
	if connection != "Upgrade" {
		return false
	}
	upgrade := c.GetHeader("Upgrade")
	return upgrade == "websocket"
}

// IsFromConsole check the request whether from Qingcloud's web console.
func IsFromConsole(c *gin.Context) bool {
	return c.GetHeader("user-agent") == "QingCloud-Web-Console"
}
