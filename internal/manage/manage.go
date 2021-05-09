package manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"time"
)

const timeLayoutStr = "2006-01-02 15:04:05"

var ctx = utils.GetRedis().Context()

func init() {
	utils.GetMysql().AutoMigrate(&User{}, &WinningInfo{})
}

/// User : struct for user
type User struct {
	ID          uint64 `gorm:"primaryKey;"`
	AccessToken string `gorm:"type:varchar(128);" json:"-"`
	TokenType   uint64 `gorm:"type:int;" json:"-"`
	Role        uint64 `gorm:"type:int;"`
	CreatedAt   time.Time
}

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
	Title       string    `gorm:"type:varchar(32); not null" json:"title"`
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
	ID          uint64  `gorm:"primary_key" json:"id"`
	Lottery     uint64  `gorm:"type:int unsigned;not null" json:"lottery"`
	Name        string  `gorm:"type:varchar(32)" json:"name"`
	Type        uint64  `gorm:"type:int" json:"type"`
	Description string  `gorm:"type:text" json:"description"`
	Pic         string  `gorm:"type:text" json:"pic"`
	Total       uint64  `gorm:"type:int" json:"total"`
	DisplayRate uint64  `gorm:"type:int" json:"displayRate"`
	Rate        uint64  `gorm:"type:int" json:"rate"`
	Value       uint64  `gorm:"type:int" json:"value"`
	Fkey        Lottery `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

//winninginfo struct
type WinningInfo struct {
	ID      uint64    `gorm:"primary_key"`
	User    uint64    `gorm:"type:int unsigned;index;not null"`
	Award   uint64    `gorm:"type:int unsigned;index;not null"`
	Lottery uint64    `gorm:"type:int unsigned;index;not null"`
	Address string    `gorm:"type:tinytext"`
	Handout bool      `gorm:"type:tinyint(1);default:0"`
	Fkey1   User      `gorm:"foreignkey:User;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fkey2   AwardInfo `gorm:"foreignkey:Award;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fkey3   Lottery   `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

//struct for high val
type Award struct {
	Award   uint64    `gorm:"primary_key; type:int unsigned"`
	Lottery uint64    `gorm:"type:int unsigned; index;not null"`
	Remain  uint64    `gorm:"type:int"`
	Fkey1   Lottery   `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fkey2   AwardInfo `gorm:"foreignkey:Award;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
		Id         uint64      `form:"id"` //lotteryid
		AwardInfos []AwardInfo `form:"awards" json:"awards"`
	}{}
	if c.ShouldBindJSON(&awards) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	tx := utils.GetMysql().Begin()
	for i, _ := range awards.AwardInfos {
		awards.AwardInfos[i].Lottery = awards.Id
		err := tx.Table("award_infos").Create(&awards.AwardInfos[i]).Error
		//low value award,set the rest of award in redis
		if awards.AwardInfos[i].Value == 0 {
			utils.GetRedis().Set(ctx, fmt.Sprintf("awards:%d", awards.AwardInfos[i].ID), awards.AwardInfos[i].Total, 0)
		} else {
			if tx.Table("awards").Create(Award{
				Award:   awards.AwardInfos[i].ID,
				Lottery: awards.AwardInfos[i].Lottery,
				Remain:  awards.AwardInfos[i].Total,
			}).Error != nil {
				tx.Rollback()
				c.Status(http.StatusBadRequest)
				return
			}
		}
		if err != nil {
			tx.Rollback()
			c.Status(http.StatusBadRequest)
			return
		}
	}
	tx.Commit()
	utils.GetRedis().Del(ctx, fmt.Sprintf("rate:%d", awards.Id))
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
	utils.GetRedis().Del(ctx, fmt.Sprintf("rate:%d", requestData.LotteryID))
	utils.GetRedis().Del(ctx, fmt.Sprintf("awards:%d", requestData.Id))
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func updateaward(c *gin.Context) {
	award := AwardInfo{}
	if c.ShouldBindJSON(&award) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err := utils.GetMysql().Model(&award).Updates(map[string]interface{}{"rate": award.Rate}).Error
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	utils.GetRedis().Del(ctx, fmt.Sprintf("rate:%d", award.Lottery))
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
	defer rows.Close()
	if err2 != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	type user struct {
		User uint64 `json:"user"`
		Name string `json:"award"`
	}
	users := make([]user, 0, requestData.Rows)
	var num uint64 = 0
	type responseData struct {
		Id     uint64 `json:"id"`
		Result []user `json:"result"`
		Page   uint64 `json:"page"`
		Rows   uint64 `json:"rows"`
		Total  uint64 `json:"total"`
	}
	for rows.Next() {
		u := user{}
		// ScanRows 将一行扫描至 user
		if utils.GetMysql().ScanRows(rows, &u) != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		users = append(users, u)
		num++
	}
	c.JSON(http.StatusOK, responseData{
		Id:     requestData.Id,
		Result: users,
		Page:   requestData.Page,
		Rows:   num,
		Total:  uint64(count),
	})
}
