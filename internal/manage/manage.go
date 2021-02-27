package manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"time"
)

const timeLayoutStr = "2006-01-02 15:04:05"

func init() {
	fmt.Print("make table....")
	utils.GetMysql().AutoMigrate(&Lottery{}, &AwardInfo{})
}

//setup the management router
func SetupManageRouter(group *gin.RouterGroup) {
	{
		group.POST("/addlottery", addlottery)
		group.POST("/updatelottery", updatelottery)
		group.POST("/addawards", addawards)
		group.DELETE("/deleteaward", deleteaward)
		group.POST("/updateaward", updateaward)
	}
}

//lottery info struct
type Lottery struct {
	ID          uint64    `gorm:"primary_key"`
	Title       string    `gorm:"type:varchar(32)" form:"title"`
	Description string    `gorm:"type:text" form:"description"`
	Permanent   uint64    `gorm:"type:int" form:"permanent"`
	Temporary   uint64    `gorm:"type:int" form:"temporary"`
	StartTime   time.Time `form:"startTime"`
	EndTime     time.Time
}

type lotteryReceived struct {
	ID          uint64 `form:"id"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Permanent   uint64 `form:"permanent"`
	Temporary   uint64 `form:"temporary"`
	StartTime   string `form:"startTime"`
	EndTime     string `form:"endTime"`
}

//awardinfo struct
type AwardInfo struct {
	ID          uint64 `gorm:"primary_key" form:"id"`
	Lottery     uint64 `gorm:"type:int" form:"lottery"`
	Name        string `gorm:"type:varchar(32)" form:"name"`
	Type        uint64 `gorm:"type:int" form:"type"`
	Description string `gorm:"type:text" form:"description"`
	Pic         string `gorm:"type:text" form:"pic"`
	Total       uint64 `gorm:"type:int" form:"total"`
	DisplayRate uint64 `gorm:"type:int" form:"displayRate"`
	Rate        uint64 `gorm:"type:int" form:"rate"`
	Value       uint64 `gorm:"type:int" form:"value"`
}

func addlottery(c *gin.Context) {
	lotteryReceived := lotteryReceived{}
	if c.ShouldBind(&lotteryReceived) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	lottery := Lottery{
		Title:       lotteryReceived.Title,
		Description: lotteryReceived.Description,
		Permanent:   lotteryReceived.Permanent,
		Temporary:   lotteryReceived.Temporary,
	}
	//transform string to Time
	loc, _ := time.LoadLocation("Local")
	if t, err := time.ParseInLocation(timeLayoutStr, lotteryReceived.StartTime, loc); err == nil {
		lottery.StartTime = t
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	if t, err := time.ParseInLocation(timeLayoutStr, lotteryReceived.EndTime, loc); err == nil {
		lottery.EndTime = t
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	model := utils.GetMysql().Table("lotteries")
	model.Create(&lottery)
	c.JSON(http.StatusOK, gin.H{"msg": "发布成功"})
}

func updatelottery(c *gin.Context) {
	lotteryReceived := lotteryReceived{}
	if c.ShouldBind(&lotteryReceived) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	lottery := Lottery{
		ID:          lotteryReceived.ID,
		Title:       lotteryReceived.Title,
		Description: lotteryReceived.Description,
		Permanent:   lotteryReceived.Permanent,
		Temporary:   lotteryReceived.Temporary,
	}
	loc, _ := time.LoadLocation("Local")
	if t, err := time.ParseInLocation(timeLayoutStr, lotteryReceived.StartTime, loc); err == nil {
		lottery.StartTime = t
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	if t, err := time.ParseInLocation(timeLayoutStr, lotteryReceived.EndTime, loc); err == nil {
		lottery.EndTime = t
	} else {
		c.Status(http.StatusBadRequest)
		return
	}
	utils.GetMysql().Model(&lottery).Updates(&lottery)
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

func addawards(c *gin.Context) {
	awards := struct {
		Id         uint64      `form:"id"`
		AwardInfos []AwardInfo `form:"awards" json:"awards"`
	}{}
	if c.ShouldBind(&awards) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	for _, award := range awards.AwardInfos {
		utils.GetMysql().Table("award_infos").Create(&AwardInfo{
			Lottery:     awards.Id,
			Name:        award.Name,
			Type:        award.Type,
			Description: award.Description,
			Pic:         award.Pic,
			Total:       award.Total,
			DisplayRate: award.DisplayRate,
			Rate:        award.Rate,
			Value:       award.Value,
		})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

func deleteaward(c *gin.Context) {
	id := struct {
		Id uint64 `form:"id" json:"id"`
	}{}
	if c.ShouldBind(&id) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	utils.GetMysql().Delete(&AwardInfo{}, id.Id)
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func updateaward(c *gin.Context) {
	award := AwardInfo{}
	if c.ShouldBind(&award) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	utils.GetMysql().Model(&award).Updates(&award)
	fmt.Println(award)
}
