package info

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLotteryInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"page":1,"rows":3}`, http.StatusOK},
		{`{"page":2,"rows":1}`, http.StatusOK},
		{`{"page":3,"rows":2}`, http.StatusNotFound},
		{`{"page":0,"rows":5}`, http.StatusNotFound},
		{`{"page":1,"rows":0}`, http.StatusNotFound},
		{`{"page":999,"rows":999}`, http.StatusNotFound},
		{`{"page":1}`, http.StatusNotFound},
		{`{}`, http.StatusNotFound},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/lottery_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestAwardsInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"lottery_id":1,"page":1,"rows":2}`, http.StatusOK},
		{`{"lottery_id":1,"page":2,"rows":1}`, http.StatusOK},
		{`{"lottery_id":1,"page":1,"rows":20}`, http.StatusOK},
		{`{"lottery_id":999,"page":1,"rows":20}`, http.StatusNotFound},
		{`{"page":1,"rows":20}`, http.StatusNotFound},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/awards_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestWinInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"user_id":26567004,"page":1,"rows":2}`, http.StatusOK},
		{`{"user_id":26567004,"page":1,"rows":10}`, http.StatusOK},
		{`{"user_id":26567004,"page":9,"rows":99}`, http.StatusNotFound},
		{`{"user_id":26567004,"page":1}`, http.StatusNotFound},
		{`{"user_id":12345678,"page":1,"rows":2}`, http.StatusNotFound},
		{`{"page":1,"rows":2}`, http.StatusNotFound},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/win_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestDrawTimes(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"user_id":26567004,"lottery_id":1}`, http.StatusOK},
		{`{"user_id":26567004,"lottery_id":2}`, http.StatusOK},
		{`{"user_id":26567004,"lottery_id":9999}`, http.StatusNotFound},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/draw_times", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}
