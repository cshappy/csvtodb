package api_csvtosqlite

import (
	svr_csvtosqlite "core/app/service/csvtosqlite"
	"core/library/response"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	// "github.com/gogf/gf/g"
	// "github.com/gogf/gf/g/net/ghttp"
)

type CsvtosqliteController struct{}

func init() {
	// 统一设置路由为小写
	g.Server().SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
}

func (d CsvtosqliteController) Csvtosqlite(r *ghttp.Request) {
	if _, err := svr_csvtosqlite.Csvtosqlite(r); err != nil {
		response.Json(r, 400, err.Error())
	}
}
