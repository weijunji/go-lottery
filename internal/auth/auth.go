package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"gorm.io/gorm/clause"
)

// User : struct for user
type User struct {
	ID          uint64 `gorm:"primaryKey;"`
	AccessToken string `gorm:"type:varchar(128);" json:"-"`
	TokenType   uint64 `gorm:"type:int;" json:"-"`
	Role        uint64 `gorm:"type:int;"`
}

// User role
const (
	RoleAdmin  = 0
	RoleNormal = 1
)

// Oauth type
const (
	OauthGithub = 0
)

// Profile : User's profile get from github
type Profile struct {
	Username string `json:"username"`
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
}

// SetupRouter : set up auth router
func SetupRouter(anonymousGroup *gin.RouterGroup, authGroup *gin.RouterGroup) {
	utils.GetMysql().AutoMigrate(&User{})
	{
		anonymousGroup.GET("/login", login)
		anonymousGroup.GET("/callback", callback)
	}
	{
		authGroup.GET("/profile", getProfile)
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
	profile := getUserFromGithub(token)
	// insert into database
	db := utils.GetMysql()
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token"}),
	}).Create(&User{profile.ID, token, OauthGithub, RoleNormal})
	// get profile from db
	user := User{ID: profile.ID}
	if err := db.First(&user).Error; err != nil {
		fmt.Printf("%+v", err)
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

func getUserFromGithub(token string) Profile {
	url := "https://api.github.com/user"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")
	res, _ := client.Do(req)
	b, _ := ioutil.ReadAll(res.Body)

	type GithubProfile struct {
		Username string `json:"name"`
		ID       uint64 `json:"id"`
		Email    string `json:"email"`
	}

	var gProfile GithubProfile
	json.Unmarshal(b, &gProfile)
	return Profile{gProfile.Username, gProfile.ID, gProfile.Email}
}

// get user's profile
func getProfile(c *gin.Context) {

}

// Login : redirect to github oauth page
func login(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/login/oauth/authorize?scope=user&client_id="+utils.GetConfig("oauth")["client_id"].(string))
}
