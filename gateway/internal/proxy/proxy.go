package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Handler(target string) gin.HandlerFunc {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	return func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	}
}
