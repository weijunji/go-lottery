package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/internal/manage"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "listening port")
	flag.Parse()
	r := setupRouter()
	r.Run(":" + strconv.Itoa(*port))
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	manage.SetupManageRouter(manageGroup)
	return r
}
