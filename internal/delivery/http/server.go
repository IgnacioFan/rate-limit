package http

import (
	"fmt"
	"go-rate-limiter/internal/service/base"
	"go-rate-limiter/internal/service/conn/redis"
	"go-rate-limiter/internal/service/ratelimiter"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	*gin.Engine
	Ratelimiter ratelimiter.Ratelimiter
	RedisClient redis.RedisClient
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		Engine:      gin.Default(),
		RedisClient: redis.NewRedisClient(),
	}
	server.SetRatelimiter()
	server.SetRouter()
	return server
}

func (s *HttpServer) SetRatelimiter() {
	s.Ratelimiter = ratelimiter.NewRatelimiter(s.RedisClient)
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
		ipAddr := strings.Split(c.ClientIP(), ":")[0]
		fmt.Println("Reveal client IP", ipAddr)
		fmt.Println(s.RedisClient.Ping(context))
		permit, remain, err := s.Ratelimiter.AcquireByIP(context, ipAddr)
		if err != nil {
			context.WithFields(logrus.Fields{"err": err, "ip": ipAddr}).Error("limiter.AcquireByIP failed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to acquire IP"})
			c.Abort()
			return
		}

		fmt.Println(context, ipAddr, permit, remain)
		if !permit {
			setAllowOrigin(c)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request"})
			c.Abort()
			return
		}
		c.Set("reqCount", remain)
		c.Next()
	}
}
