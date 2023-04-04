package http

import "github.com/gin-gonic/gin"

func setAllowOrigin(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	c.Header("Access-Control-Allow-Origin", origin)
}
