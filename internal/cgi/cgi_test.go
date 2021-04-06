package cgi

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	SetupRouter(gin.Default())
}

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func TestForward(t *testing.T) {
	assert := assert.New(t)
	r := gin.Default()
	r.Any("/", Forward("http://www.baidu.com"))
	w := CreateTestResponseRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
}
