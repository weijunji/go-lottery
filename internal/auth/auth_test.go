package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/middleware"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var normalToken, _ = utils.GenerateToken(99998888, RoleNormal, time.Minute*10)
var adminToken, _ = utils.GenerateToken(99998887, RoleAdmin, time.Minute*10)
var superToken, _ = utils.GenerateToken(99998886, RoleSuperAdmin, time.Minute*10)

func setup() {
	gin.SetMode(gin.TestMode)
	db := utils.GetMysql()
	db.Create(&User{ID: 99999999, AccessToken: "test_test_test", TokenType: OauthGithub, Role: RoleNormal})
	db.Create(&User{ID: 99999998, AccessToken: "test_test_test", TokenType: OauthGithub, Role: RoleNormal})
	db.Create(&User{ID: 99999997, AccessToken: "test_test_test", TokenType: OauthGithub, Role: RoleAdmin})
	db.Create(&User{ID: 99999996, AccessToken: "test_test_test", TokenType: OauthGithub, Role: RoleNormal})
	db.Create(&User{ID: 99999995, AccessToken: "test_test_test", TokenType: OauthGithub, Role: RoleAdmin})
}

func teardown() {
	db := utils.GetMysql()
	db.Delete(&User{}, []uint64{99999999, 99999998, 99999997, 99999996, 99999995})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestUpdateUser(t *testing.T) {
	tg := User{ID: 99999999}
	utils.GetMysql().First(&tg)
	assert.Equal(t, tg.Role, RoleNormal, "user's role should be normal")

	r := gin.Default()
	g := r.Group("/auth", middleware.AuthMiddleware())
	authGroup := g.Group("/", middleware.LoginRequired())
	SetupRouter(g, authGroup)

	cases := []struct {
		token  string
		body   string
		expect int
	}{
		{normalToken, "", 403},
		{adminToken, "", 403},
		{superToken, "", 400},
		{superToken, "{\"id\": 99999999, \"role\": 0}", 200},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/auth/user", strings.NewReader(c.body))
		req.Header.Add("Authorization", c.token)
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}

	tg = User{ID: 99999999}
	utils.GetMysql().First(&tg)
	assert.Equal(t, tg.Role, RoleAdmin, "user's role should be changed")
}

func TestDeleteUser(t *testing.T) {
	r := gin.Default()
	g := r.Group("/auth", middleware.AuthMiddleware())
	authGroup := g.Group("/", middleware.LoginRequired())
	SetupRouter(g, authGroup)

	cases := []struct {
		token  string
		body   string
		expect int
	}{
		{normalToken, "", 403},
		{superToken, "", 400},
		{adminToken, "{\"id\": 99999998}", 200},
		{adminToken, "{\"id\": 99999997}", 404},
		{superToken, "{\"id\": 99999996}", 200},
		{superToken, "{\"id\": 99999995}", 200},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/auth/user", strings.NewReader(c.body))
		req.Header.Add("Authorization", c.token)
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}

	tg := User{ID: 99999997}
	rows := utils.GetMysql().First(&tg).RowsAffected
	assert.Equal(t, rows, int64(1))
	assert.Equal(t, tg.Role, RoleAdmin, "user 99999997 should not be delete")

	rows = utils.GetMysql().Find(&User{}, []uint64{99999998, 99999996, 99999995}).RowsAffected
	assert.Equal(t, rows, int64(0), "users should be delete")
}
