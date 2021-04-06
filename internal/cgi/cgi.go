package cgi

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var targetMap = map[string]string{
	"Auth": "http://localhost:8082",
	"Lottery": "http://localhost:8083",
	"Info": "http://localhost:8084",
	"Manage": "http://localhost:8085",
}

func SetupRouter(r *gin.Engine) {
	r.Any("/auth", Forward(targetMap["Auth"]))
	r.Any("/lottery", Forward(targetMap["Lottery"]))
	r.Any("/info", Forward(targetMap["Info"]))
	r.Any("/manage", Forward(targetMap["Manage"]))
}

func Forward(target string) func(c *gin.Context) {
	return func(c *gin.Context) {
		HostReverseProxy(c.Writer, c.Request, target)
	}
}

func HostReverseProxy(w http.ResponseWriter, req *http.Request, target string) {
	remote, err := url.Parse(target)
	if err != nil {
		log.Errorf("err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
}
