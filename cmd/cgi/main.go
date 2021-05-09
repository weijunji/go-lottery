package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/internal/cgi"
	"strconv"
)

func main(){
	port := flag.Int("port", 8081, "listening port")
	flag.Parse()
	r := setupRouter()
	r.Run(":" + strconv.Itoa(*port))
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	cgi.SetupRouter(r)
	return r
}
