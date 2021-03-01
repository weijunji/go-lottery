package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/weijunji/go-lottery/pkgs/middleware"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"gorm.io/gorm/clause"
)

// User : struct for user
type User struct {
	ID          uint64 `gorm:"primaryKey;"`
	AccessToken string `gorm:"type:varchar(128);" json:"-"`
	TokenType   uint64 `gorm:"type:int;" json:"-"`
	Role        uint64 `gorm:"type:int;"`
	CreatedAt   time.Time
}

// User role
const (
	RoleAdmin      uint64 = 0
	RoleNormal     uint64 = 1
	RoleSuperAdmin uint64 = 2
)

// Oauth type
const (
	OauthGithub uint64 = 0
)

// Profile : User's profile get from github
type Profile struct {
	Username string `json:"username"`
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
}

// SetupRouter : set up auth router
func SetupRouter(anonymousGroup *gin.RouterGroup, authGroup *gin.RouterGroup) {
	if utils.GetMysql().AutoMigrate(&User{}) != nil {
		log.Fatal("Auto migrate failed")
	}
	{
		anonymousGroup.GET("/login", login)
		anonymousGroup.GET("/callback", callback)
	}
	{
		authGroup.GET("/profile", getProfile)
		authGroup.PUT("/user", middleware.SuperAdminRequired(), updateUser)
		authGroup.DELETE("/user", middleware.AdminRequired(), deleteUser)
	}
}

func callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	// get access token from github
	config := utils.GetConfig("oauth")
	body := fmt.Sprintf("{\"client_id\":\"%s\", \"client_secret\": \"%s\", \"code\": \"%s\"}", config["client_id"].(string), config["client_secret"].(string), code)
	resp, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBufferString(body))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	token := strings.Split(strings.Split(string(b), "&")[0], "=")[1]
	// get user profile from github
	profile, err := getUserFromGithub(token)
	if err != nil {
		log.Error("get user profile failed")
		c.Status(http.StatusInternalServerError)
		return
	}
	// insert into database
	db := utils.GetMysql()
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token"}),
	}).Create(&User{ID: profile.ID, AccessToken: token, TokenType: OauthGithub, Role: RoleNormal})
	// get profile from db
	user := User{ID: profile.ID}
	if err := db.First(&user).Error; err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}
	// generate jwt key
	jwt, err := utils.GenerateToken(user.ID, user.Role, 3*24*time.Hour)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, gin.H{"token": jwt, "role": user.Role, "profile": profile})
}

var ErrorWrongToken = errors.New("wrong token")

func getUserFromGithub(token string) (Profile, error) {
	url := "https://api.github.com/user"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")
	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		return Profile{}, ErrorWrongToken
	}
	b, _ := ioutil.ReadAll(res.Body)

	type GithubProfile struct {
		Username string `json:"name"`
		ID       uint64 `json:"id"`
		Email    string `json:"email"`
	}

	var gProfile GithubProfile
	if err := json.Unmarshal(b, &gProfile); err != nil {
		log.Fatal("Json unmarshal failed")
	}
	return Profile{gProfile.Username, gProfile.ID, gProfile.Email}, nil
}

// get user's profile
func getProfile(c *gin.Context) {
	userInfo, ok := c.Get("userinfo")
	if !ok {
		log.Fatal("Get user info failed: user should login")
	}
	id := userInfo.(middleware.Userinfo).ID
	db := utils.GetMysql()

	user := User{ID: id}
	if err := db.First(&user).Error; err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}

	if profile, err := getUserFromGithub(user.AccessToken); err == nil {
		c.JSON(200, profile)
	} else {
		c.Status(http.StatusUnauthorized)
	}
}

// Login : redirect to github oauth page
func login(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/login/oauth/authorize?scope=user&client_id="+utils.GetConfig("oauth")["client_id"].(string))
}

func updateUser(c *gin.Context) {
	request := struct {
		ID   uint64 `json:"id"`
		Role uint64 `json:"role"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	db := utils.GetMysql()
	if db.Model(&User{}).Where("id = ?", request.ID).Update("role", request.Role).RowsAffected == 0 {
		c.Status(http.StatusNotFound)
	} else {
		c.Status(http.StatusOK)
	}
}

func deleteUser(c *gin.Context) {
	request := struct {
		ID uint64 `json:"id"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	info, _ := c.Get("userinfo")
	db := utils.GetMysql()

	if info.(middleware.Userinfo).Role == RoleAdmin {
		// admin can delete normal
		if db.Where("role = ?", RoleNormal).Delete(&User{}, request.ID).RowsAffected == 0 {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusOK)
		}
	} else {
		// super can delete all except super
		if db.Where("role != ?", RoleSuperAdmin).Delete(&User{}, request.ID).RowsAffected == 0 {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusOK)
		}
	}
}
