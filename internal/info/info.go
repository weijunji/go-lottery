package info

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/weijunji/go-lottery/pkgs/utils"
	myProto "github.com/weijunji/go-lottery/proto"
)

//const timeFormat = "2006-01-02 15:04:05"
var (
	rdb = utils.GetRedis()
)

func init() {
	_ = utils.GetMysql().AutoMigrate(&Users{}, &Lotteries{}, &AwardInfos{}, &Awards{}, &WinningInfos{})
}

//load routers
func LoadRouter(r *gin.RouterGroup) {
	{
		r.GET("/lottery_info", LotteryInfo)
		r.GET("/awards_info", AwardsInfo)
		r.GET("/win_info", WinInfo)
		r.GET("/draw_times", DrawTimes)
	}
}

// Users [...]
type Users struct {
	ID          uint64    `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`
	AccessToken string    `gorm:"column:access_token;type:varchar(128)" json:"-"`
	TokenType   uint64    `gorm:"column:token_type;type:bigint" json:"-"`
	Role        uint64    `gorm:"column:role;type:bigint" json:"role"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3)" json:"created_at"`
}
// TableName get sql table name.
func (m *Users) TableName() string {
	return "users"
}

//Lotteries: struct for lotteries
type Lotteries struct {
	ID          uint64    `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"lottery_id"`
	Title       string    `gorm:"column:title;type:varchar(32);not null" json:"lottery_title"`
	Description string    `gorm:"column:description;type:text" json:"lottery_description"`
	Permanent   uint64    `gorm:"column:permanent;type:bigint" json:"permanent"`
	Temporary   uint64    `gorm:"column:temporary;type:bigint" json:"temporary"`
	StartTime   time.Time `gorm:"column:start_time;type:datetime(3)" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time;type:datetime(3)" json:"end_time"`
}
func (m *Lotteries) TableName() string {
	return "lotteries"
}

//AwardInfos: struct for award_infos
type AwardInfos struct {
	ID          uint64    `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`
	Lottery     uint64    `gorm:"column:lottery;type:bigint unsigned;not null" json:"lottery"`
	Fkey   		Lotteries `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string    `gorm:"column:name;type:varchar(32)" json:"name"`
	Type        int64     `gorm:"column:type;type:bigint" json:"type"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Pic         string    `gorm:"column:pic;type:text" json:"pic"`
	Total       uint64    `gorm:"column:total;type:bigint" json:"total"`
	DisplayRate uint64    `gorm:"column:display_rate;type:bigint" json:"display_rate"`
	Rate        uint64    `gorm:"column:rate;type:bigint" json:"rate"`
	Value       uint64    `gorm:"column:value;type:bigint" json:"value"`
}
func (m *AwardInfos) TableName() string {
	return "award_infos"
}

// Awards [...]
type Awards struct {
	Award      uint64     `gorm:"primaryKey;column:award;type:bigint unsigned;not null" json:"-"`
	Fkey2 	   AwardInfos `gorm:"foreignkey:Award;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Lottery    uint64     `gorm:"index:idx_awards_lottery;column:lottery;type:bigint unsigned;not null" json:"lottery"`
	Fkey1  	   Lotteries  `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reamin     int64      `gorm:"column:reamin;type:bigint" json:"reamin"`
}
// TableName get sql table name.
func (m *Awards) TableName() string {
	return "awards"
}

//WinningInfo: struct for winning_infos
type WinningInfos struct {
	ID         uint64     `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`
	User       uint64     `gorm:"index:idx_winning_infos_user;column:user;type:bigint unsigned;not null" json:"user"`
	Fkey1      Users      `gorm:"foreignkey:User;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Award      uint64     `gorm:"index:idx_winning_infos_award;column:award;type:bigint unsigned;not null" json:"award"`
	Fkey2 	   AwardInfos `gorm:"foreignkey:Award;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Lottery    uint64     `gorm:"index:idx_winning_infos_lottery;column:lottery;type:bigint unsigned;not null" json:"lottery"`
	Fkey3      Lotteries  `gorm:"foreignkey:Lottery;association_foreignkey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Address    string     `gorm:"column:address;type:tinytext" json:"address"`
	Handout    bool       `gorm:"column:handout;type:tinyint(1);default:0" json:"handout"`
}
// TableName get sql table name.
func (m *WinningInfos) TableName() string {
	return "winning_infos"
}

//LotteryInfoRes: struct for LotteryInfo()
type LotteryInfoRes struct {
	LotteryItems []Lotteries `json:"lottery_items"`
	Page         uint64      `json:"page"`
	Rows         uint64      `json:"rows"`
	Total        int64       `json:"total"`
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
	page := request.Page
	rows := request.Rows
	db := utils.GetMysql()

	var lotteries []Lotteries
	var lotteryCount int64
	if err := db.Table("lotteries").Where("end_time > NOW()").Find(&lotteries).Count(&lotteryCount).Error; err != nil {
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

type AwardItem struct {
	ID          uint64 `json:"award_id"`
	Name        string `json:"award_name"`
	Description string `json:"award_description"`
	Pic         string `json:"pic"`
	Total       uint64 `json:"total"`
	DisplayRate uint64 `json:"display_rate"`
	Value       uint64 `json:"value"`
}

//AwardInfoRes: struct for AwardsInfo()
type AwardInfoRes struct {
	LotteryId          uint64      `json:"lottery_id"`
	LotteryTitle       string      `json:"lottery_title"`
	LotteryDescription string      `json:"lottery_description"`
	Awards             []AwardItem `json:"awards"`
	Page               uint64      `json:"page"`
	Rows               uint64      `json:"rows"`
	Total              int64       `json:"total"`
}

//Get prize information for lottery
func AwardsInfo(c *gin.Context) {
	request := struct {
		LotteryId uint64 `json:"lottery_id"`
		Page      uint64 `json:"page"`
		Rows      uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	lotteryId := request.LotteryId
	page := request.Page
	rows := request.Rows
	db := utils.GetMysql()

	var awardsCount int64
	var awardsInfos []AwardInfos
	if err := db.Table("award_infos").Where("lottery = ?", lotteryId).Find(&awardsInfos).Count(&awardsCount).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	flag := judgePageRange(page, rows, uint64(awardsCount))
	if flag == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var lotteryTemp Lotteries
	if err := db.Table("lotteries").Where("id = ?", lotteryId).Find(&lotteryTemp).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	awardsInfos = awardsInfos[(page-1)*rows : flag]

	var awardsRequest []AwardItem
	var t AwardItem
	for _, v := range awardsInfos {
		t.ID = v.ID
		t.Name = v.Name
		t.Description = v.Description
		t.Pic = v.Pic
		t.Total = v.Total
		t.DisplayRate = v.DisplayRate
		t.Value = v.Value
		awardsRequest = append(awardsRequest, t)
	}
	res := AwardInfoRes{
		lotteryId,
		lotteryTemp.Title,
		lotteryTemp.Description,
		awardsRequest,
		page,
		rows,
		awardsCount,
	}
	c.JSON(http.StatusOK, res)
}

type WinItem struct {
	Lottery uint64 `json:"lottery_id"`
	Title   string `json:"lottery_title"`
	AwardId uint64 `json:"award_id"`
	Name    string `json:"award_name"`
	Address string `json:"address"`
	Handout bool   `json:"handout"`
}
type WinningInfoRes struct {
	UserId uint64    `json:"user_id"`
	Awards []WinItem `json:"awards"`
	Page   uint64    `json:"page"`
	Rows   uint64    `json:"rows"`
	Total  int64     `json:"total"`
}

//Get the user winning information of the lottery
func WinInfo(c *gin.Context) {
	request := struct {
		UserId uint64 `json:"user_id"`
		Page   uint64 `json:"page"`
		Rows   uint64 `json:"rows"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userId := request.UserId
	page := request.Page
	rows := request.Rows

	var winningCount int64
	var winRes []WinItem

	db := utils.GetMysql()
	err := db.Table("winning_infos").
		Select("winning_infos.lottery, lotteries.title, winning_infos.award, award_infos.name, winning_infos.address, winning_infos.handout").
		Joins("INNER JOIN lotteries ON winning_infos.lottery=lotteries.id").
		Joins("INNER JOIN award_infos ON winning_infos.award=award_infos.id").
		Where("winning_infos.user = ?", userId).
		Find(&winRes).Count(&winningCount).Error
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	flag := judgePageRange(page, rows, uint64(winningCount))
	if flag == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	winReq := winRes[(page-1)*rows : flag]
	res := WinningInfoRes{
		userId,
		winReq,
		page,
		rows,
		winningCount,
	}
	c.JSON(http.StatusOK, res)
}

type returnType struct {
	Permanent uint64 `json:"permanent"`
	Temporary uint64 `json:"temporary"`
}

//Query user's remaining lottery draws
func DrawTimes(c *gin.Context) {
	request := struct {
		UserId    uint64 `json:"user_id"`
		LotteryId uint64 `json:"lottery_id"`
	}{}
	if c.ShouldBindJSON(&request) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userId := request.UserId
	lotteryId := request.LotteryId

	ctx := context.Background()
	if rdb == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var toReturn returnType
	ans, err := rdb.Get(ctx, "remain:"+strconv.Itoa(int(lotteryId))+":"+strconv.Itoa(int(userId))).Result()
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
		toReturn.Permanent = lottery.Permanent
		toReturn.Temporary = lottery.Temporary

		timesLeft := &myProto.UserTimes{}
		timesLeft.Permanent = uint32(lottery.Permanent)
		timesLeft.Temporary = uint32(lottery.Temporary)
		timesLeft.Update = ptypes.TimestampNow()
		// save to redis
		buffer, err := proto.Marshal(timesLeft)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		rdb.Set(ctx, "remain:"+strconv.FormatUint(lotteryId, 10)+":"+strconv.FormatUint(userId, 10), buffer, 0)
	} else {
		// found in redis
		data := &myProto.UserTimes{}
		if err := proto.Unmarshal([]byte(ans), data); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		toReturn.Permanent = uint64(data.Permanent)
		toReturn.Temporary = uint64(data.Temporary)

		// Update for the first time today
		timeStamp0, _ := getTimestamp()       // timestamp for 00:00:00 of today
		timeStampRedis := data.Update.Seconds // timestamp for last update
		if time.Unix(timeStampRedis, 0).Before(time.Unix(timeStamp0, 0)) {
			lottery, flag := QueryLotteryById(lotteryId)
			if flag != 0 {
				c.Status(http.StatusNotFound)
				return
			}
			// update "temporary" and "update"
			data.Temporary = uint32(lottery.Temporary)
			data.Update = ptypes.TimestampNow()
			buffer, err := proto.Marshal(data)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			rdb.Set(ctx, "remain:"+strconv.FormatUint(lotteryId, 10)+":"+strconv.FormatUint(userId, 10), buffer, 0)
			toReturn.Temporary = lottery.Temporary
		}
	}
	c.JSON(http.StatusOK, toReturn)
}

func QueryLotteryById(lotteryId uint64) (Lotteries, uint64) {
	db := utils.GetMysql()
	var lottery Lotteries
	if err := db.Where("id = ?", lotteryId).Find(&lottery).Error; err != nil {
		return lottery, 1
	}
	return lottery, 0
}

// Solve page numbering issues
func judgePageRange(page, rows, cnt uint64) uint64 {
	flag := page * rows
	// out of range
	if (page != 1 && flag > cnt) || flag == 0 {
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
