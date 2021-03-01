package manage

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"strings"
	"time"
)

const timeLayoutStr = "2006-01-02 15:04:05"

//var ctx = utils.GetRedis().Context()

//setup the management router
func SetupManageRouter(group *gin.RouterGroup) {
	{
		group.POST("/addlottery", addlottery)
		group.POST("/updatelottery", updatelottery)
		group.POST("/addawards", addawards)
		group.DELETE("/deleteaward", deleteaward)
		group.POST("/updateaward", updateaward)
		group.GET("/getawardinfolist", getawardinfolist)
	}
}

//lottery info struct
type Lottery struct {
	ID          uint64    `gorm:"primary_key"`
	Title       string    `gorm:"type:varchar(32)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Permanent   uint64    `gorm:"type:int" json:"permanent"`
	Temporary   uint64    `gorm:"type:int" json:"temporary"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time
}

type lotteryReceived struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Permanent   uint64 `json:"permanent"`
	Temporary   uint64 `json:"temporary"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

//awardinfo struct
type AwardInfo struct {
	ID          uint64 `gorm:"primary_key" json:"id"`
	Lottery     uint64 `gorm:"type:int" json:"lottery"`
	Name        string `gorm:"type:varchar(32)" json:"name"`
	Type        uint64 `gorm:"type:int" json:"type"`
	Description string `gorm:"type:text" json:"description"`
	Pic         string `gorm:"type:text" json:"pic"`
	Total       uint64 `gorm:"type:int" json:"total"`
	DisplayRate uint64 `gorm:"type:int" json:"displayRate"`
	Rate        uint64 `gorm:"type:int" json:"rate"`
	Value       uint64 `gorm:"type:int" json:"value"`
}

type WinningInfo struct {
	ID      uint64 `gorm:"primary_key"`
	User    uint64 `gorm:"type:int"`
	Award   uint64 `gorm:"type:int"`
	Lottery uint64 `gorm:"type:int"`
	Address string `gorm:"type:tinytext"`
	Handout bool   `gorm:"type:tinyint(1)"`
}

func addlottery(c *gin.Context) {
	lotteryReceived := lotteryReceived{}
	if c.ShouldBindJSON(&lotteryReceived) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	lottery := Lottery{
		Title:       lotteryReceived.Title,
		Description: lotteryReceived.Description,
		Permanent:   lotteryReceived.Permanent,
		Temporary:   lotteryReceived.Temporary,
	}
	fmt.Println(lotteryReceived)
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
	err := utils.GetMysql().Table("lotteries").Create(&lottery).Error
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "发布成功"})
}

func updatelottery(c *gin.Context) {
	lotteryReceived := lotteryReceived{}
	if c.ShouldBindJSON(&lotteryReceived) != nil {
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
	if utils.GetMysql().Model(&lottery).Updates(&lottery).Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

func addawards(c *gin.Context) {
	awards := struct {
		Id         uint64      `form:"id"`
		AwardInfos []AwardInfo `form:"awards" json:"awards"`
	}{}
	if c.ShouldBindJSON(&awards) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	tx := utils.GetMysql().Table("award_infos").Begin()
	for _, award := range awards.AwardInfos {
		err := tx.Create(&AwardInfo{
			Lottery:     awards.Id,
			Name:        award.Name,
			Type:        award.Type,
			Description: award.Description,
			Pic:         award.Pic,
			Total:       award.Total,
			DisplayRate: award.DisplayRate,
			Rate:        award.Rate,
			Value:       award.Value,
		}).Error
		if err != nil {
			tx.Rollback()
			c.Status(http.StatusBadRequest)
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

func deleteaward(c *gin.Context) {
	requestData := struct {
		Id        uint64 `form:"id" json:"award"`
		LotteryID uint64 `json:"lottery"`
	}{}
	if c.ShouldBindJSON(&requestData) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if utils.GetMysql().Delete(&AwardInfo{}, requestData.Id).RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	//Invalidate the probability of prizes and the number of low-value prizes in redis
	//utils.GetRedis().Del(ctx,fmt.Sprintf("rate:%d",requestData.Id)).Err()
	//utils.GetRedis().SRem(ctx,fmt.Sprintf("lottery:%d",requestData.LotteryID),requestData.LotteryID).Err()
	//utils.GetRedis().Del(ctx,fmt.Sprintf("awards:%d",requestData.Id))
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func updateaward(c *gin.Context) {
	award := AwardInfo{}
	if c.ShouldBindJSON(&award) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err := utils.GetMysql().Model(&award).Updates(&award).Error
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//Invalidate the probability of prizes and the number of low-value prizes in redis
	//utils.GetRedis().Del(ctx,fmt.Sprintf("rate:%d",award.ID)).Err()
	//utils.GetRedis().Del(ctx,fmt.Sprintf("awards:%d",award.ID))
	c.Status(http.StatusOK)
}

func getawardinfolist(c *gin.Context) {
	//struct for requestData
	requestData := struct {
		Id   uint64 `json:"id"`
		Page uint64 `json:"page"`
		Rows uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&requestData) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//the number of total records
	var count int64
	err1 := utils.GetMysql().Model(&WinningInfo{}).Where("lottery=?", requestData.Id).Count(&count).Error
	if err1 != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	rows, err2 := utils.GetMysql().Model(&AwardInfo{}).Raw("SELECT t1.user,t2.name from (SELECT user, award  FROM winning_infos WHERE lottery = ?) as t1 inner join award_infos as t2 on t1.award=t2.id order by t2.value desc limit ?,?", requestData.Id, (requestData.Page-1)*requestData.Rows, requestData.Rows).Rows()
	if rows != nil {
		defer rows.Close().Error()
	}
	if err2 != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user := struct {
		User uint64 `json:"user"`
		Name string `json:"award"`
	}{}

	var responseData strings.Builder
	responseData.WriteString(fmt.Sprintf(`{"id":%d,"result:[`, requestData.Id))

	num := 0
	for rows.Next() {
		num++
		if num != 1 {
			responseData.WriteString(",")
		}
		// ScanRows 将一行扫描至 user
		if utils.GetMysql().ScanRows(rows, &user) != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		temp, _ := json.Marshal(user)
		responseData.Write(temp)
	}
	responseData.WriteString(fmt.Sprintf(`],"page":%d,"rows":%d,"total":%d}`, requestData.Page, num, count))
	c.String(http.StatusOK, responseData.String())
}
