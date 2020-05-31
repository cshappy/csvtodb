package response

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)
// 标准返回结果数据结构封装
// 返回固定数据结构的JSON
// status: 错误码(200: 成功，403:失败，>403:其他错误码)
// info: 请求结果信息
// data: 请求结果，根据不同接口返回结果的数据结构不同

func Json(r *ghttp.Request, err int, msg string, data ...interface{})  {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"status": err,
		"info" : msg,
		"data" : responseData,
	})
}

func Json1(r *ghttp.Request, err int, msg string, data ...interface{})  {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(responseData)
}

