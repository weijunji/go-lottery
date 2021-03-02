package manage

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddLottery(t *testing.T) {
	r := gin.Default()
	manageGroup := r.Group("/manage")
	SetupManageRouter(manageGroup)
	testData := []struct {
		body string
		code int
	}{
		{`{"title": "嘉年华plus", "description": "好礼相送", "permanent":10, "temporary":6, "startTime":"2020-02-27 15:04:05", "endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
		{`{"title": "嘉年华plus2","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-2-2 15:04:05","endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
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
		{`{"id":2, "title": "嘉年华plus", "description": "好礼相送", "permanent":1, "temporary":6, "startTime":"2020-02-27 15:04:05", "endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
		{`{"id":2, "title": "嘉年华plus2","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-2-2 15:04:05","endTime":"2021-3-27 00:00:00"}`, http.StatusBadRequest},
		{`{"id":33, "title": "嘉年华plus2","description": "好礼相送","permanent":10,"temporary":6,"startTime":"2020-02-02 15:04:05","endTime":"2021-03-27 00:00:00"}`, http.StatusOK},
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
		{`{"id":19,
				"awards": [
				   {
						"name":"iphone",
						"type":1,
						"description":"手机",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1,
						"displayrate":20000,
						"rate":10000,
						"value":30000
				   }
    			]
		}`, http.StatusOK},
		{`{"id":19,
				"awards": [
				   {
						"name":"ipad",
						"type":1,
						"description":"平板",
						"pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
						"total":1,
						"displayrate":20000,
						"rate":20000,
						"value":5000
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
						"value":10000
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
		{`{"award": 23, "lottery":19}`, http.StatusOK},
		{`{"award": 23, "lottery":19}`, http.StatusNotFound},
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
    		"id": 22,     
			"lottery": 19,
    		"name":"iphone",
            "type":1,
            "description":"手机",
            "pic":"https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=2&pn=1&spn=0&di=7700&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3363295869%2C2467511306&os=892371676%2C71334739&simid=4203536407%2C592943110&adpicid=0&lpn=0&ln=30&fr=ala&fm=&sme=&cg=&bdtype=0&oriquery=%E5%9B%BE%E7%89%87&objurl=https%3A%2F%2Fgimg2.baidu.com%2Fimage_search%2Fsrc%3Dhttp%3A%2F%2Fa0.att.hudong.com%2F30%2F29%2F01300000201438121627296084016.jpg%26refer%3Dhttp%3A%2F%2Fa0.att.hudong.com%26app%3D2002%26size%3Df9999%2C10000%26q%3Da80%26n%3D0%26g%3D0n%26fmt%3Djpeg%3Fsec%3D1617003197%26t%3Dc44d15e73e5f0050501f06230d7f4091&fromurl=ippr_z2C%24qAzdH3FAzdH3Fooo_z%26e3Bfhyvg8_z%26e3Bv54AzdH3F4AzdH3Fetjo_z%26e3Brir%3Fwt1%3Dmb9l9&gsm=1&islist=&querylist=",
            "total":2,
            "displayrate":20000,
            "rate":10000,
            "value":10000
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
	}
	for _, data := range testData {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/manage/getawardinfolist", strings.NewReader(data.body))
		r.ServeHTTP(w, req)
		assert.Equal(t, data.code, w.Code)
	}
}
