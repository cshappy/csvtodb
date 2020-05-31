package router

import (
	api_csvtosqlite "core/app/api/csvtosqlite"
	api_history "core/app/api/history"

	"github.com/gogf/gf/frame/g"
)

func init() {
	server := g.Server()
	server.BindObject("/api/history", new(api_history.HistoryController))
	server.BindObject("/api/csvtosqlite", new(api_csvtosqlite.CsvtosqliteController))
}
