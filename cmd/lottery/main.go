package lottery

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/internal/lottery"
	"github.com/weijunji/go-lottery/pkgs/middleware"
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
	g := r.Group("/lottery", middleware.AuthMiddleware())
	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	authGroup := g.Group("/", middleware.LoginRequired())
	lottery.SetupRouter(g, authGroup)
	return r
}
