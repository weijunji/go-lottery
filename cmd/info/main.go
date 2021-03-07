package main

import (
	"flag"
	"github.com/weijunji/go-lottery/internal/info"
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
	infoGroup := r.Group("/info")
	info.LoadRouter(infoGroup)
	return r
}