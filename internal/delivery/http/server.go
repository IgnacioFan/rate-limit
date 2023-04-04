package http

import (
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/ratelimiter"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*gin.Engine
	Ratelimiter ratelimiter.Ratelimiter
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		Engine:      gin.Default(),
		Ratelimiter: ratelimiter.NewRatelimiter(),
	}
	server.SetRouter()
	return server
}

func (s *HttpServer) SetRouter() {
	v1 := s.Group("/api/v1")
	v1.Use(
		s.AcquireIP(),
	)
	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong!")
	})
}

func (s *HttpServer) AcquireIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		context := base.Background()
		clientIP := strings.Split(c.ClientIP(), ":")[0]
		fmt.Println("Reveal client IP", clientIP)

		permit, count := s.Ratelimiter.AcquireByIP(context, clientIP)

		fmt.Println(context, clientIP, permit, count)
		if !permit {
			setAllowOrigin(c)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request"})
			c.Abort()
			return
		}
		c.Set("reqCount", count)
		c.Next()
	}
}
