package main

import (
	"flag"
	"github.com/weijunji/go-lottery/internal/info"
	"github.com/weijunji/go-lottery/pkgs/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("port", 8080, "listening port")
	flag.Parse()
	r := setupInfoRouter()
	r.Run(":" + strconv.Itoa(*port))
}

func setupInfoRouter() *gin.Engine {
	r := gin.Default()
	g := r.Group("/info", middleware.AuthMiddleware())
	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	infoGroup := g.Group("/", middleware.LoginRequired())
	info.LoadRouter(g, infoGroup)
	return r
}