package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setup() {
	gin.SetMode(gin.TestMode)

	db, _ := utils.GetMysql().DB()
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1000, 'temptest001', 1, 2, NOW(), '2029-1-1')")
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1001, 'temptest002', 5, 5, '1999-1-1', '2020-12-31')")
	db.Exec("INSERT INTO lotteries(id, title, permanent, temporary, start_time, end_time) VALUES (1002, 'temptest003', 5, 5, '2029-1-1', '2029-1-2')")
}

func teardown() {
	db, _ := utils.GetMysql().DB()
	db.Exec("DELETE FROM lotteries WHERE id IN (1000, 1001, 1002)")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestAddLottery(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"title": "嘉年华plus1", "description": "好礼相送", "permanent":10, "temporary":6, "startTime":"2020-02-27 15:04:05", "endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
		{`{"title": "嘉年华plus2", "description": "好礼相送", "permanent":10, "temporary":6, "startTime":"2020-2-2 15:04:05", "endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
		{`{"title": "嘉年华plus3", "description": "好礼相送", "permanent":10, "temporary":6, "startTime":"2020-02-02 15:04:05", "endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/manage/addlottery", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}

func TestUpdateLottery(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"id":1000, "title": "temptest001", "description": "好礼相送", "permanent":1, "temporary":6, "startTime":"2020-02-27 15:04:05", "endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
		{`{"id":1001, "title": "嘉年华plus1","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-2-2 15:04:05","endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
		{`{"id":1002, "title": "嘉年华plus1","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-02-02 15:04:05","endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
		{`{"id":33, "title": "嘉年华plus1","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-02-02 15:04:05","endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/manage/updatelottery", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}

func TestAddawards(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"id":1000,
				"awards": [
				   {
						"name":"iphone",
						"type":1,
						"description":"手机",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1,
						"displayrate":20000,
						"rate":10000,
						"value":1
				   }
    			]
		}`, http.StatusOK},
		{`{"id":1000,
				"awards": [
				   {
						"name":"ipad",
						"type":1,
						"description":"平板",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1,
						"displayrate":20000,
						"rate":8000,
						"value":1
				   }
    			]
		}`, http.StatusOK},
		{`{"id":1000,
				"awards": [
				   {
						"name":"再来一次",
						"type":1,
						"description":"平板",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1000,
						"displayrate":20000,
						"rate":200000,
						"value":0
				   }
    			]
		}`, http.StatusOK},
		{`{"id":33,
				"awards": [
				   {
						"name":"iphone",
						"type":1,
						"description":"手机",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1,
						"displayrate":20000,
						"rate":20000,
						"value":1
				   }
    			]
		}`, http.StatusBadRequest},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/manage/addawards", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}

func TestDeleteaward(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"award": 1, "lottery":1}`, http.StatusOK},
		{`{"award": 1, "lottery":1}`, http.StatusNotFound},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/manage/deleteaward", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}
func TestUpdateaward(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{   
    		"id": 2,     
			"lottery": 20,
    		"name":"ipad",
			"type":1,
			"description":"平板",
			"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
			"total":1,
			"displayrate":20000,
			"rate":80090,
			"value":25000
			}`, http.StatusOK},
		{`{   
    		"id": 77, 
			"lottery": 19,
    		"name":"iphone",
            "type":1,
            "description":"手机",
            "pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
            "total":2,
            "displayrate":20000,
            "rate":20000,
            "value":10000
			}`, http.StatusOK},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/manage/updateaward", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}
func TestGetawardinfolist(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"id":1,"page":1,"rows":10}`, http.StatusOK},
		{`{"id":1,"page":2,"rows":10}`, http.StatusOK},
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/manage/getawardinfolist", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}
