package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/auth"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var normalToken, _ = utils.GenerateToken(99998888, auth.RoleNormal, time.Minute*10)
var adminToken, _ = utils.GenerateToken(99998887, auth.RoleAdmin, time.Minute*10)
var superToken, _ = utils.GenerateToken(99998886, auth.RoleSuperAdmin, time.Minute*10)

func TestAuthMiddleware(t *testing.T) {
	r := gin.Default()
	g := r.Group("/", AuthMiddleware())
	authGroup := g.Group("/", LoginRequired())
	f := func(c *gin.Context) {
		c.Status(200)
	}
	{
		g.GET("/anonymous", f)
	}
	{
		authGroup.GET("/login_required", f)
		authGroup.GET("/admin_required", AdminRequired(), f)
		authGroup.GET("/super_required", SuperAdminRequired(), f)
	}

	cases := []struct {
		token  string
		target string
		expect int
	}{
		{"", "/anonymous", 200},
		{normalToken, "/anonymous", 200},
		{adminToken, "/anonymous", 200},
		{superToken, "/anonymous", 200},

		{"", "/login_required", 401},
		{normalToken, "/login_required", 200},
		{adminToken, "/login_required", 200},
		{superToken, "/login_required", 200},

		{"", "/admin_required", 401},
		{normalToken, "/admin_required", 403},
		{adminToken, "/admin_required", 200},
		{superToken, "/admin_required", 200},

		{"", "/super_required", 401},
		{normalToken, "/super_required", 403},
		{adminToken, "/super_required", 403},
		{superToken, "/super_required", 200},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", c.target, nil)
		req.Header.Add("Authorization", c.token)
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}
