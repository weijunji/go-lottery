package manage

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/weijunji/go-lottery/pkgs/utils"
	myproto "github.com/weijunji/go-lottery/proto"
	"net/http"
	"strings"
	"time"
)

const timeLayoutStr = "2006-01-02 15:04:05"

var ctx = utils.GetRedis().Context()

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

//winninginfo struct
type WinningInfo struct {
	ID      uint64 `gorm:"primary_key"`
	User    uint64 `gorm:"type:int"`
	Award   uint64 `gorm:"type:int"`
	Lottery uint64 `gorm:"type:int"`
	Address string `gorm:"type:tinytext"`
	Handout bool   `gorm:"type:tinyint(1)"`
}

//struct for high val
type Award struct {
	Award   uint64 `gorm:"primary_key; type:int"`
	Lottery uint64 `gorm:"type:int"`
	Reamin  uint64 `gorm:"type:int"`
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
	rate := &myproto.LotteryRates{
		Total: 0,
		Rates: []*myproto.LotteryRates_AwardRate{},
	}
	rateProto, err := proto.Marshal(rate)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	utils.GetRedis().Set(ctx, fmt.Sprintf("rate:%d", lottery.ID), string(rateProto), 0)
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
	tx := utils.GetMysql().Table("award_infos").Begin()
	for i, _ := range awards.AwardInfos {
		awards.AwardInfos[i].Lottery = awards.Id
		err := tx.Create(&awards.AwardInfos[i]).Error
		if err != nil {
			tx.Rollback()
			c.Status(http.StatusBadRequest)
			return
		}
	}
	tx.Commit()
	val := utils.GetRedis().Get(ctx, fmt.Sprintf("rate:%d", awards.Id)).Val()
	rate := &myproto.LotteryRates{}
	if proto.Unmarshal([]byte(val), rate) != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for _, award := range awards.AwardInfos {
		rate.Total++
		rate.Rates = append(rate.Rates, &myproto.LotteryRates_AwardRate{
			Id:   award.ID,
			Rate: uint32(award.Rate),
		})
		//low value award,set the rest of award in redis
		if award.Value < 20000 {
			utils.GetRedis().Set(ctx, fmt.Sprintf("awards:%d", award.ID), award.Total, 0)
		} else {
			utils.GetMysql().Table("awards").Create(Award{
				Award:   award.ID,
				Lottery: award.Lottery,
				Reamin:  award.Total,
			})
		}
	}
	rateProto, _ := proto.Marshal(rate)
	utils.GetRedis().Set(ctx, fmt.Sprintf("rate:%d", awards.Id), string(rateProto), 0)
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
	val := utils.GetRedis().Get(ctx, fmt.Sprintf("rate:%d", requestData.LotteryID)).Val()
	rate := &myproto.LotteryRates{}
	if proto.Unmarshal([]byte(val), rate) != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for i, v := range rate.Rates {
		if v.Id == requestData.Id {
			rate.Total--
			rate.Rates = append(rate.Rates[:i], rate.Rates[i+1:]...)
		}
	}
	rateProto, _ := proto.Marshal(rate)
	utils.GetRedis().Set(ctx, fmt.Sprintf("rate:%d", requestData.LotteryID), string(rateProto), 0)
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
	rate := &myproto.LotteryRates{}
	val := utils.GetRedis().Get(ctx, fmt.Sprintf("rate:%d", award.Lottery)).Val()
	if proto.Unmarshal([]byte(val), rate) != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for i, _ := range rate.Rates {
		if rate.Rates[i].Id == award.ID {
			rate.Rates[i].Rate = uint32(award.Rate)
		}
	}
	rateProto, _ := proto.Marshal(rate)
	utils.GetRedis().Set(ctx, fmt.Sprintf("rate:%d", award.Lottery), string(rateProto), 0)
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
