package log

import (
	"io"
	"os"
	"time"

	"github.com/gogf/gf/net/ghttp"
)

func init() {
	start()
	s := ghttp.GetServer()
	s.SetLogPath("log")
	s.SetAccessLogEnabled(true)
	// s.SetLogHandler(Writelog)
}

func Writelog(r *ghttp.Request, errs ...interface{}) {
	timeStr := time.Now().String()[:23]
	uri := r.RequestURI
	var msg string

	filePath := "log/" + timeStr[:10]
	if len(errs) > 0 {
		msg = "time: " + timeStr + "	remote: " + r.RemoteAddr + "	interface: " + uri + "	error: " + errs[0].(error).Error() + "\n"
	}
	msg = "time: " + timeStr + "	remote: " + r.RemoteAddr + "	interface: " + uri + "	event: " + "visit" + "\n"
	writeFile(filePath, []byte(msg), 0644)
}

func start() {
	timeStr := time.Now().String()[:23]
	var msg = "time: " + timeStr + "	event: " + "server is start" + "\n"
	filePath := "log/" + timeStr[:10]
	writeFile(filePath, []byte(msg), 0644)
}
func writeFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	return err
}
