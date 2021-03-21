package info

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/weijunji/go-lottery/pkgs/auth"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	gin.SetMode(gin.TestMode)
	db := utils.GetMysql()
	_ = db.AutoMigrate(&Users{}, &Lotteries{}, &AwardInfos{}, &Awards{}, &WinningInfos{})

	timeStr := []string {
		"2020-02-27 15:04:05.000",
		"2021-03-27 00:00:00.000",
		"2021-03-04 13:59:47.000",
		"2021-03-31 13:59:55.000",
		"2021-03-05 17:46:44.000",
		"2021-04-02 17:46:52.000",
		"2021-03-01 17:47:24.000",
		"2021-03-04 17:47:30.000",
	}
	var res []time.Time
	for i := range timeStr {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr[i], time.Local)
		res = append(res, t)
	}

	db.Create(&Users{ID: 26567004, AccessToken: "test_test_test", TokenType: 0, Role: auth.RoleNormal})

	db.Create(&Lotteries{ID: 1, Title: "plus1", Description: "description1", Permanent: 1, Temporary: 6, StartTime: res[0], EndTime: res[1]})
	db.Create(&Lotteries{ID: 2, Title: "plus2", Description: "description2", Permanent: 2, Temporary: 6, StartTime: res[2], EndTime: res[3]})
	db.Create(&Lotteries{ID: 3, Title: "plus3", Description: "description3", Permanent: 3, Temporary: 10, StartTime: res[4], EndTime: res[5]})
	db.Create(&Lotteries{ID: 4, Title: "plus4", Description: "description4", Permanent: 1, Temporary: 1, StartTime: res[6], EndTime: res[7]})

	db.Create(&AwardInfos{ID: 1, Lottery: 1, Name: "ipad", Type: 1, Description: "ipad", Pic: "test", Total: 1, DisplayRate: 20000, Rate: 80090, Value: 25000})
	db.Create(&AwardInfos{ID: 2, Lottery: 1, Name: "again", Type: 1, Description: "again", Pic: "test", Total: 1000, DisplayRate: 20000, Rate: 200000, Value: 100})

	db.Create(&Awards{Award: 1, Lottery: 1, Remain: 1})

	db.Create(&WinningInfos{ID: 1, User: 26567004, Award: 1, Lottery: 1, Address: "test", Handout: false})
	db.Create(&WinningInfos{ID: 2, User: 26567004, Award: 2, Lottery: 1, Address: "test", Handout: true})
}

func teardown() {
	db := utils.GetMysql()
	db.Delete(&Users{}, []uint64{26567004})
	db.Delete(&Lotteries{}, []uint64{1, 2, 3, 4})
	db.Delete(&AwardInfos{}, []uint64{1, 2})
	db.Delete(&Awards{}, []uint64{1})
	db.Delete(&WinningInfos{}, []uint64{1, 2})
}


func TestLotteryInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"page":1,"rows":3}`, http.StatusOK},
		{`{"page":2,"rows":1}`, http.StatusOK},
		{`{"page":999,"rows":999}`, http.StatusOK},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/lottery_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestAwardsInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"lottery_id":1,"page":1,"rows":2}`, http.StatusOK},
		{`{"lottery_id":1,"page":2,"rows":1}`, http.StatusOK},
		{`{"lottery_id":1,"page":1,"rows":20}`, http.StatusOK},
		{`{"lottery_id":999,"page":1,"rows":20}`, http.StatusOK},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/awards_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestWinInfo(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"user_id":26567004,"page":1,"rows":2}`, http.StatusOK},
		{`{"user_id":26567004,"page":1,"rows":10}`, http.StatusOK},
		{`{"user_id":26567004,"page":9,"rows":99}`, http.StatusOK},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/win_info", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}

func TestDrawTimes(t *testing.T) {
	r := gin.Default()
	infoGroup := r.Group("/info")
	LoadRouter(infoGroup)

	cases := []struct {
		body string
		expect int
	}{
		{`{"user_id":26567004,"lottery_id":1}`, http.StatusOK},
		{`{"user_id":26567004,"lottery_id":2}`, http.StatusOK},
		{`{"user_id":26567004,"lottery_id":9999}`, http.StatusNotFound},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info/draw_times", strings.NewReader(c.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, c.expect, w.Code)
	}
}
