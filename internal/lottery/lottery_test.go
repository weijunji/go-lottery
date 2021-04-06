package lottery

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/weijunji/go-lottery/pkgs/auth"
	"github.com/weijunji/go-lottery/pkgs/middleware"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	pb "github.com/weijunji/go-lottery/proto"
)

func setup() {
	gin.SetMode(gin.TestMode)
	rate := pb.LotteryRates{
		Total: 100,
		Rates: [] *pb.LotteryRates_AwardRate {
			{Id: 1, Rate: 20, Value: pb.LotteryRates_LOW_VAL},
			{Id: 2, Rate: 30, Value: pb.LotteryRates_LOW_VAL},
			{Id: 3, Rate: 50, Value: pb.LotteryRates_HIGH_VAL},
		},
	}
	rds := utils.GetRedis()
	bytes, _ := proto.Marshal(&rate)
	rds.Set(context.Background(), "rate:99999999", bytes, 0)

	db, _ := utils.GetMysql().DB()
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1000, 'temptest001', 1, 2, NOW(), '2029-1-1')")
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1001, 'temptest002', 5, 5, '1999-1-1', '2020-12-31')")
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1002, 'temptest003', 5, 5, '2029-1-1', '2029-1-2')")

	t, _ := proto.Marshal(&pb.UserTimes{Permanent: 0, Temporary: 3, Update: timestamppb.New(time.Unix(0, 0))})
	rds.Set(context.Background(), "remain:1000:629", t, 0)

	rds.Set(context.Background(), "awards:628629", 2, 0)

	db.Exec("INSERT INTO award_infos(id, lottery, rate, value) VALUES (999, 1001, 200000, 1)")
	db.Exec("INSERT INTO awards(award, lottery, remain) VALUES (999, 1001, 2)")

	db.Exec("INSERT INTO award_infos(id, lottery, rate, value) VALUES (1000, 1000, 200000, 1)")
	db.Exec("INSERT INTO award_infos(id, lottery, rate, value) VALUES (1001, 1000, 300000, 1)")
	db.Exec("INSERT INTO award_infos(id, lottery, rate, value) VALUES (1002, 1000, 500000, 0)")
	db.Exec("INSERT INTO awards(award, lottery, remain) VALUES (1000, 1000, 2)")
	db.Exec("INSERT INTO awards(award, lottery, remain) VALUES (1001, 1000, 5)")
	rds.Set(context.Background(), "awards:1002", 100, 0)
}

func teardown() {
	rds := utils.GetRedis()
	rds.Del(context.Background(), "rate:99999999")

	db, _ := utils.GetMysql().DB()
	db.Exec("DELETE FROM lotteries WHERE id IN (1000, 1001, 1002)")
	rds.Del(context.Background(), "lottery_duration:1000", "lottery_duration:1001", "lottery_duration:1001")
	rds.Del(context.Background(), "lottery_times:1000", "lottery_times:1001", "lottery_times:1001")
	rds.Del(context.Background(), "remain:1000:628", "remain:1000:629")
	rds.Del(context.Background(), "awards:628629")
	db.Exec("DELETE FROM award_infos WHERE id IN (999, 1000, 1001, 1002)")
	db.Exec("DELETE FROM awards WHERE id IN (999, 1000, 1001)")
	rds.Del(context.Background(), "rate:1000")
	rds.Del(context.Background(), "awards:1002")
	rds.Del(context.Background(), "remain:1000:99998888")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetRate(t *testing.T) {
	assert := assert.New(t)
	rates := GetRate(context.Background(), 99999999)
	assert.NotEqual(rates, nil)
	assert.Equal(rates.GetTotal(), uint32(100))
	assert.Equal(rates.GetRates()[2].Rate, uint32(50))

	rates = GetRate(context.Background(), 99999998)
	assert.Nil(rates)

	rates = GetRate(context.Background(), 1000)
	assert.Equal(rates.GetTotal(), uint32(1000000))
}

func TestProcessLottery(t *testing.T) {
	assert := assert.New(t)
	const times = 1000
	rate := pb.LotteryRates{
		Total: 500000,
		Rates: [] *pb.LotteryRates_AwardRate {
			{Id: 1, Rate: 400000},
			{Id: 2, Rate: 100000},
		},
	}
	var win [3]uint32
	for i := 0; i < times; i++ {
		award := processLottery(&rate)
		if award == nil {
			win[0]++
		} else {
			win[award.GetId()]++
		}
	}
	assert.Greater(win[0], uint32(400))
	assert.Greater(win[1], uint32(350))
	assert.Greater(win[2], uint32(70))
	assert.Less(win[0], uint32(600))
	assert.Less(win[1], uint32(450))
	assert.Less(win[2], uint32(130))
}


func TestGetDuration(t *testing.T) {
	assert := assert.New(t)
	// from sql
	res := GetLotteryDuration(context.Background(), 1000)
	assert.NotNil(res)
	now := time.Now()
	assert.True(now.After(res.GetStart().AsTime()))
	assert.True(now.Before(res.GetEnd().AsTime()))

	// from redis
	res = GetLotteryDuration(context.Background(), 1000)
	assert.NotNil(res)

	res = GetLotteryDuration(context.Background(), 999)
	assert.Nil(res)
}

func TestGetTimes(t *testing.T) {
	assert := assert.New(t)
	// from sql
	res := GetLotteryTimes(context.Background(), 1000)
	assert.NotNil(res)
	assert.Equal(uint32(1), res.Permanent)
	assert.Equal(uint32(2), res.Temporary)

	// from redis
	res = GetLotteryTimes(context.Background(), 1000)
	assert.NotNil(res)
}

func TestDecreaseTimes(t *testing.T) {
	assert := assert.New(t)

	assert.True(decreaseTimes(context.Background(), 1000, 628))
	assert.True(decreaseTimes(context.Background(), 1000, 628))
	assert.True(decreaseTimes(context.Background(), 1000, 628))
	assert.False(decreaseTimes(context.Background(), 1000, 628))

	assert.True(decreaseTimes(context.Background(), 1000, 629))
	assert.True(decreaseTimes(context.Background(), 1000, 629))
	assert.False(decreaseTimes(context.Background(), 1000, 629))
}

func TestDecrLowValue(t *testing.T) {
	assert := assert.New(t)

	assert.True(decreaseAwardLowVal(context.Background(), 628629))
	assert.True(decreaseAwardLowVal(context.Background(), 628629))
	assert.False(decreaseAwardLowVal(context.Background(), 628629))
}

func TestDecrHighValue(t *testing.T) {
	assert := assert.New(t)

	assert.True(decreaseAwardHighVal(context.Background(), 999))
	assert.True(decreaseAwardHighVal(context.Background(), 999))
	assert.False(decreaseAwardHighVal(context.Background(), 999))
}

func TestLottery(t *testing.T) {
	assert := assert.New(t)
	r := gin.Default()
	g := r.Group("/lottery", middleware.AuthMiddleware())
	authGroup := g.Group("/", middleware.LoginRequired())
	SetupRouter(g, authGroup)

	token, _ := utils.GenerateToken(99998888, auth.RoleNormal, time.Minute*10)

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/lottery/once?id=1000", nil)
		req.Header.Add("Authorization", token)
		r.ServeHTTP(w, req)
		assert.Equal(200, w.Code)
		assert.Contains(w.Body.String(), "true")
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lottery/once?id=1000", nil)
	req.Header.Add("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Contains(w.Body.String(), "false")
}

func TestLotteryError(t *testing.T) {
	assert := assert.New(t)
	r := gin.Default()
	g := r.Group("/lottery", middleware.AuthMiddleware())
	authGroup := g.Group("/", middleware.LoginRequired())
	SetupRouter(g, authGroup)

	token, _ := utils.GenerateToken(99998888, auth.RoleNormal, time.Minute*10)

	cases := []struct {
		url string
		expect int
	}{
		{"/lottery/once", http.StatusBadRequest},
		{"/lottery/once?id=abs", http.StatusBadRequest},
		{"/lottery/once?id=999", http.StatusNotFound},
		{"/lottery/once?id=1001", http.StatusNotFound},
		{"/lottery/once?id=1002", http.StatusNotFound},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", c.url, nil)
		req.Header.Add("Authorization", token)
		r.ServeHTTP(w, req)
		assert.Equal(c.expect, w.Code)
	}
}
