package info

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"strconv"
	"time"
)

//const timeFormat = "2006-01-02 15:04:05"
var (
	rdb = utils.GetRedis()
)

//load routers
func LoadRouter(r *gin.RouterGroup) {
	{
		r.GET("/lottery_info", 	LotteryInfo)
		r.GET("/awards_info",  	AwardsInfo)
		r.GET("/win_info", 		WinInfo)
		r.GET("/draw_times", 	DrawTimes)
	}
}

//Lottery: struct for lotteries
type Lottery struct {
	ID			uint64		`json:"lottery_id"`
	Title		string		`json:"lottery_title"`
	Description string		`json:"lottery_description"`
	Permanent	uint64		`json:"permanent"`
	Temporary 	uint64		`json:"temporary"`
	StartTime 	time.Time	`json:"start_time"`
	EndTime		time.Time	`json:"end_time"`
}
//LotteryInfoRes: struct for LotteryInfo()
type LotteryInfoRes struct {
	LotteryItems 	[]Lottery	`json:"lottery_items"`
	Page 			uint64 		`json:"page"`
	Rows 			uint64		`json:"rows"`
	Total			int64		`json:"total"`
}
//Get a list of ongoing lottery
func LotteryInfo(c *gin.Context) {
	request := struct {
		Page uint64 `json:"page"`
		Rows uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	db := utils.GetMysql()

	var lotteries 			[]Lottery
	var lotteryCount		int64
	page := request.Page
	rows := request.Rows

	if err := db.AutoMigrate(&Lottery{}); err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}
	if err := db.Where("end_time > NOW()").Find(&lotteries).Count(&lotteryCount).Error; err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}
	//out of range
	flag := judgePageRange(page, rows, uint64(lotteryCount))
	if flag == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	lotteriesRequest := lotteries[(page-1)*rows : flag]
	res := LotteryInfoRes{
		lotteriesRequest,
		page,
		rows,
		lotteryCount,
	}
	c.JSON(http.StatusOK, res)
}


//AwardInfo: struct for award_infos
type AwardInfo struct {
	ID				uint64	`json:"award_id"`
	Name			string	`json:"award_name"`
	Description		string	`json:"award_description"`
	Pic 			string	`json:"pic"`
	Total 			uint64	`json:"total"`
	DisplayRate		uint64	`json:"display_rate"`
	Value 			uint64	`json:"value"`
}
//AwardInfoRes: struct for AwardsInfo()
type AwardInfoRes struct {
	LotteryId			uint64		`json:"lottery_id"`
	LotteryTitle		string		`json:"lottery_title"`
	LotteryDescription	string		`json:"lottery_description"`
	Awards				[]AwardInfo	`json:"awards"`
	Page 				uint64 		`json:"page"`
	Rows 				uint64		`json:"rows"`
	Total				int64		`json:"total"`
}
//Get prize information for lottery
func AwardsInfo(c *gin.Context) {
	request := struct {
		LotteryId uint64 `json:"lottery_id"`
		Page uint64 `json:"page"`
		Rows uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	db := utils.GetMysql()
	if err := db.AutoMigrate(&AwardInfo{}); err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}

	var awardsCount int64
	var awardsInfo []AwardInfo
	if err := db.Where("lottery = ?", request.LotteryId).Find(&awardsInfo).Count(&awardsCount).Error; err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}

	page := request.Page
	rows := request.Rows
	flag := judgePageRange(page, rows, uint64(awardsCount))
	if flag == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	if err := db.AutoMigrate(&Lottery{}); err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}

	var lotteryTemp Lottery
	if err := db.Table("lotteries").Where("id = ?", request.LotteryId).Find(&lotteryTemp).Error; err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}

	LotteryId 			:= request.LotteryId
	LotteryTitle 		:= lotteryTemp.Title
	LotteryDescription	:= lotteryTemp.Description
	awardsRequest := awardsInfo[(page-1)*rows : flag]
	res := AwardInfoRes{
		LotteryId,
		LotteryTitle,
		LotteryDescription,
		awardsRequest,
		page,
		rows,
		awardsCount,
	}
	c.JSON(http.StatusOK, res)
}


//WinningInfo: struct for winning_infos
type WinningInfo struct {
	ID      uint64 `gorm:"primary_key"`
	User    uint64 `gorm:"type:int"`
	Award   uint64 `gorm:"type:int"`
	Lottery uint64 `gorm:"type:int"`
	Address string `gorm:"type:tinytext"`
	Handout bool   `gorm:"type:tinyint(1)"`
}
type Award struct {
		Lottery		uint64	`json:"lottery_id"`
		Title		string	`json:"lottery_title"`
		Award		uint64	`json:"award_id"`
		Name		string	`json:"award_name"`
		Address 	string	`json:"address"`
		Handout		bool	`json:"handout"`
}
type WinningInfoRes struct {
	UserId 	uint64		`json:"user_id"`
	Awards	[]Award		`json:"awards"`
	Page 	uint64 		`json:"page"`
	Rows 	uint64		`json:"rows"`
	Total	int64		`json:"total"`
}
//Get the user winning information of the lottery
func WinInfo(c *gin.Context) {
	request := struct {
		UserId uint64 `json:"user_id"`
		Page uint64 `json:"page"`
		Rows uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userId 	:= request.UserId
	page 	:= request.Page
	rows 	:= request.Rows

	db := utils.GetMysql()
	if err := db.AutoMigrate(&WinningInfo{}); err != nil {
		log.Errorf("%+v", err)
		c.Status(http.StatusNotFound)
		return
	}
	var winningCount int64
	var awardRes []Award
	err2 := db.Table("winning_infos").
		Select("winning_infos.lottery, lotteries.title, winning_infos.award, award_infos.name, winning_infos.address, winning_infos.handout").
		Joins("INNER JOIN lotteries ON winning_infos.lottery=lotteries.id").
		Joins("INNER JOIN award_infos ON winning_infos.award=award_infos.id").
		Where("winning_infos.user = ?", userId).
		Find(&awardRes).
		Count(&winningCount).Error
	if err2 != nil {
		c.Status(http.StatusNotFound)
		return
	}

	flag := judgePageRange(page, rows, uint64(winningCount))
	if flag == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	awardReq := awardRes[(page-1)*rows : flag]
	res := WinningInfoRes{
		userId,
		awardReq,
		page,
		rows,
		winningCount,
	}
	c.JSON(http.StatusOK, res)
}


type UserTimes struct {
	Permanent	uint64	`json:"permanent"`
	Temporary	uint64	`json:"temporary"`
}
type redisFormat struct {
	Permanent	uint64		`json:"permanent"`
	Temporary	uint64		`json:"temporary"`
	Update		time.Time	`json:"update"`
}
func QueryLotteryById(lotteryId uint64) (Lottery, uint64) {
	db := utils.GetMysql()
	var lottery Lottery
	if err := db.AutoMigrate(&Lottery{}); err != nil {
		return lottery, 1
	}
	if err := db.Where("id = ?", lotteryId).Find(&lottery).Error; err != nil {
		return lottery, 1
	}
	return lottery, 0
}
//Query user's remaining lottery draws
func DrawTimes(c *gin.Context) {
	request := struct {
		UserId 		uint64 `json:"user_id"`
		LotteryId	uint64	`json:"lottery_id"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userId 		:= request.UserId
	lotteryId	:= request.LotteryId

	ctx := context.Background()
	if rdb == nil {
		log.Fatal("connect redis failed")
		c.Status(http.StatusInternalServerError)
		return
	}

	var userTimes	UserTimes
	var redisFormat redisFormat
	ans, err := rdb.Get( ctx, "remain:" + strconv.Itoa(int(lotteryId)) + ":" + strconv.Itoa(int(userId)) ).Result()
	if err != nil {
		// not found in redis, create one and insert into redis
		lottery, flag := QueryLotteryById(lotteryId)
		if flag != 0 {
			c.Status(http.StatusNotFound)
			return
		}
		//no data in lottery
		if lottery.ID == 0 {
			c.Status(http.StatusNotFound)
			return
		}
		userTimes.Permanent = lottery.Permanent
		userTimes.Temporary = lottery.Temporary

		redisFormat.Permanent = lottery.Permanent
		redisFormat.Temporary = lottery.Temporary
		redisFormat.Update	  = time.Now()
		// save to redis
		toRedis, err 	:= json.Marshal(redisFormat)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		rdb.Set(ctx, "remain:" + strconv.FormatUint(lotteryId, 10) + ":" + strconv.FormatUint(userId, 10),
			string(toRedis), 0)
	} else {
		// found in redis
		if err := json.Unmarshal([]byte(ans), &redisFormat); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		userTimes.Permanent = redisFormat.Permanent
		userTimes.Temporary = redisFormat.Temporary

		timeStamp0, _ := getTimestamp()
		// Update for the first time today
		if redisFormat.Update.Before(time.Unix(timeStamp0, 0)) {
			//???
			lottery, flag := QueryLotteryById(lotteryId)
			if flag != 0 {
				c.Status(http.StatusNotFound)
				return
			}
			// update "temporary" and "update"
			redisFormat.Temporary = lottery.Temporary
			redisFormat.Update	  = time.Now()

			redisFormatStr, err := json.Marshal(redisFormat)
			if  err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			rdb.Set(ctx, "remain:" + strconv.FormatUint(lotteryId, 10) + ":" + strconv.FormatUint(userId, 10), string(redisFormatStr), 0)
			userTimes.Temporary = lottery.Temporary
		}
	}
	c.JSON(http.StatusOK, userTimes)
}

// Solve page numbering issues
func judgePageRange(page, rows, cnt uint64) uint64 {
	flag := page * rows
	// out of range
	if (page != 1 && flag > cnt ) || flag == 0 {
		return 0
	}
	// return the maximum rows
	if flag > cnt {
		return cnt
	} else {
		return flag
	}
}

// Get the timestamp of 0 o'clock and 24 o'clock today
func getTimestamp() (beginTimeNum, endTimeNum int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	beginTimeNum = t.Unix()
	endTimeNum = beginTimeNum + 86400
	return beginTimeNum, endTimeNum
}
