package main

import (
	_ "core/router"
	_ "core/library/log"
	"github.com/gogf/gf/frame/g"
)

func main() {
	// g.Server是单利模式  任何时候返回的都是同一个对象
	s := g.Server()
	// 设置cookie的时间 0 表示不启用cookie  -1表示关闭浏览器 cookie失效
	s.SetCookieMaxAge(60 * 60 * 7)
	s.SetIndexFolder(true)
	s.SetServerRoot("public")
	s.SetPort(8090)
	s.Run()
}