package api_history

import (
	svr_history "core/app/service/history"
	"core/library/response"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	// "fmt"
	// "reflect"
)

type HistoryController struct{}

func init() {
	// 统一设置路由为小写
	g.Server().SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
}

func (d HistoryController) Gethistorian(r *ghttp.Request) {

	if code, err := svr_history.Gethistorian(r.URL.RawQuery); err != nil {
		// fmt.Println(r.URL.RawQuery)
		response.Json(r, 400, err.Error())
	} else {
		result := code.([]interface{})
		result1 := (result[0].(map[string]interface{}))
		if result1["ErrorCode"] == 0 || result1["ErrorCode"] == -14 || result1["ErrorCode"] == nil {
			response.Json1(r, 200, "", map[string]interface{}{
				"ErrorCode":    0,
				"ErrorMessage": nil,
				"Data":         code,
			})
		} else if result1["ErrorCode"] == 10 {
			response.Json1(r, 200, "", map[string]interface{}{
				"ErrorCode":    10,
				"ErrorMessage": "Invalid tagName",
			})
		}
	}
}
