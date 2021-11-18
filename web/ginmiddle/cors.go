package ginmiddle

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	cfg := cors.Config{
		AllowAllOrigins:        true,
		AllowOrigins:           nil,
		AllowOriginFunc:        nil,
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type", "Content-MD5", "Authorization", "Date", "X-Date"},
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}
	return cors.New(cfg)
}
