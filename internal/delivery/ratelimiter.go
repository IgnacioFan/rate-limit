package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*gin.Engine
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		Engine: gin.Default(),
	}
	server.SetRouter()
	return server
}

func (s *HttpServer) SetRouter() {
	v1 := s.Group("/api/v1")
	v1.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong!")
	})
}
