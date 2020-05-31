package svr_reserve

import (
	"fmt"
	// "github.com/gogf/gf/g"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
	_ "github.com/mattn/go-adodb"
	// "reflect"
)

func Gethistorian(data string) (interface{}, error) {
	// var result = make(map[string]interface{})
	var result = []interface{}{}
	// var tagname=data["tagname"]
	// var timestart=data["timestart"]
	// var timeend  =data["timeend"]
	// var ip=data["ip"]
	var tagname string
	var timestart string
	var timeend string
	var ip string
	var target string
	var timestart1 string
	var timeend1 string
	var intervalmilliseconds string
	var calculationMode string
	var numberofsamples string
	var samplingmode string
	var timezone string
	arr1 := strings.Split(data, "&")
	// var tag string
	for i := 0; i < len(arr1); i++ {
		arr2 := strings.Split(arr1[i], "=")
		if arr2[0] == "tagname" {
			tagname = arr2[1]
		} else if arr2[0] == "start" {
			timestart = arr2[1]
			timestart1 = arr2[1]
		} else if arr2[0] == "end" {
			timeend = arr2[1]
			timeend1 = arr2[1]
		} else if arr2[0] == "ip" {
			ip = arr2[1]
		} else if arr2[0] == "target" {
			target = arr2[1]
		} else if arr2[0] == "intervalMilliseconds" {
			intervalmilliseconds = arr2[1]
		} else if arr2[0] == "calculationMode" {
			calculationMode = arr2[1]
		} else if arr2[0] == "numberofsamples" {
			numberofsamples = arr2[1]
		} else if arr2[0] == "samplingmode" {
			samplingmode = arr2[1]
		} else if arr2[0] == "timezone" {
			timezone = arr2[1]
		}

	}
	if tagname == "" {
		var s1 = make(map[string]interface{})
		s1["ErrorCode"] = 10
		s1["ErrorMessage"] = "Invalid tagName"
		result = append(result, []interface{}{s1}...)
		return result, nil
	}

	db, err := sql.Open("adodb", "Provider=ihOLEDB.iHistorian.1;User ID=;Password=;Data Source="+ip+";")
	if err != nil {
		var s1 = make(map[string]interface{})
		s1["ErrorCode"] = 10
		s1["ErrorMessage"] = "Invalid sql"
		result = append(result, []interface{}{s1}...)
		return result, nil
	}
	timestart, err = getTime(timestart)
	fmt.Println(timestart)
	if err != nil {
		var s1 = make(map[string]interface{})
		s1["status"] = 405
		s1["error"] = "error"
		s1["message"] = "invlaid Starttime"
		result = append(result, []interface{}{s1}...)
		return result, nil
	}
	timeend, err = getTime(timeend)
	fmt.Println(timeend)
	if err != nil {
		var s1 = make(map[string]interface{})
		s1["status"] = 405
		s1["error"] = "error"
		s1["message"] = "invlaid Endtime"
		result = append(result, []interface{}{s1}...)
		return result, nil
	}
	arr := strings.Split(tagname, ";")
	var tag string
	for i := 1; i < len(arr); i++ {
		tag = tag + " or tagname=" + arr[i]
	}
	var sqll string
	if len(timestart1) == 0 || (len(timeend1) == 0 && len(timestart1) == 0) {
		sqll = "SELECT tagname,quality,value,timestamp FROM ihRawData where tagname=" + arr[0] + tag
	} else if len(timeend1) == 0 {
		sqll = "SELECT tagname,quality,value,timestamp FROM ihRawData where tagname=" + arr[0] + tag + " and timestamp>='" + timestart
	} else {
		sqll = "SELECT tagname,quality,value,timestamp FROM ihRawData where tagname=" + arr[0] + tag + " and timestamp>='" + timestart + "' and timestamp<='" + timeend + "'"
	}
	if len(intervalmilliseconds) != 0 {
		sqll = sqll + " and intervalmilliseconds=" + intervalmilliseconds
	}

	if len(calculationMode) != 0 {
		sqll = sqll + " and calculationMode=" + calculationMode
	}

	if len(numberofsamples) != 0 {
		sqll = sqll + " and numberofsamples=" + numberofsamples
	}

	if len(samplingmode) != 0 {
		sqll = sqll + " and samplingmode=" + samplingmode
	}

	if len(timezone) != 0 {
		sqll = sqll + " and timezone=" + timezone
	}

	rows, err := db.Query(sqll)
	fmt.Println(sqll)
	if err != nil {
		fmt.Println("出错")
		panic(err)
		fmt.Println(err)
	}
	if target == "" {
		var i string
		var i1 = 1
		var samples = []interface{}{}
		var s1 = make(map[string]interface{})
		var sample1 = make(map[string]interface{})
		for rows.Next() {
			var tagname *ole.VARIANT
			var quality *ole.VARIANT
			var value *ole.VARIANT
			var timestamp string
			rows.Scan(&tagname, &quality, &value, &timestamp)
			if i1 == 1 {
				i = tagname.ToString()
			}
			if i == tagname.ToString() {
				s1["TagName"] = tagname.ToString()
				s1["ErrorCode"] = 0
				sample1["Value"] = value.Value()
				sample1["Quality"] = quality.Value()
				sample1["Timestamp"] = timestamp
				samples = append(samples, []interface{}{sample1}...)
				sample1 = make(map[string]interface{})
			} else {
				s1["Samples"] = samples
				result = append(result, []interface{}{s1}...)
				samples = []interface{}{}
				s1 = make(map[string]interface{})
				i = tagname.ToString()
				s1["TagName"] = tagname.ToString()
				s1["ErrorCode"] = 0
				sample1["Value"] = value.Value()
				sample1["Quality"] = quality.Value()
				sample1["Timestamp"] = timestamp
				samples = append(samples, []interface{}{sample1}...)
				sample1 = make(map[string]interface{})
			}
			i1 = i1 + 1
		}
		s1["Samples"] = samples
		if len(samples) == 0 {
			s1["ErrorCode"] = -14
		}
		result = append(result, []interface{}{s1}...)
		return result, nil
	} else if target == "echart" {
		var s2 = make(map[string]interface{})
		var i2 string
		var i3 = 1
		var data1 = []interface{}{}
		var data = make(map[string]interface{})
		var data2 = []interface{}{}
		// var data=make(map[string]interface{})
		for rows.Next() {
			var tagname *ole.VARIANT
			var quality *ole.VARIANT
			var value *ole.VARIANT
			var timestamp string
			// var data1=""
			s2 = make(map[string]interface{})
			rows.Scan(&tagname, &quality, &value, &timestamp)
			if i3 == 1 {
				i2 = tagname.ToString()
			}
			if i2 == tagname.ToString() {
				s2["name"] = tagname.ToString()
				data1 = append(data1, timestamp)
				number := stringzero(strconv.FormatFloat(value.Value().(float64), 'f', 4, 64))
				data1 = append(data1, number)
				data["value"] = data1
				data2 = append(data2, data)
				data1 = []interface{}{}
				data = make(map[string]interface{})
			} else {
				s2["name"] = i2
				i2 = tagname.ToString()
				s2["data"] = data2
				result = append(result, s2)
				// fmt.Println(result)
				data2 = []interface{}{}
				s2 = make(map[string]interface{})
				data1 = []interface{}{}
				data = make(map[string]interface{})
				data1 = append(data1, timestamp)
				number := stringzero(strconv.FormatFloat(value.Value().(float64), 'f', 4, 64))
				data1 = append(data1, number)
				data["value"] = data1
				data2 = append(data2, data)
				data1 = []interface{}{}
			}
			i3 = i3 + 1
			s2["data"] = data2
		}
		result = append(result, s2)
		return result, nil
	}
	return result, nil

}
func stringzero(num string) string {
	for i := 0; i < 4; i++ {
		if num[len(num)-1:len(num)] == "0" {
			num = num[:len(num)-1]
		}
	}
	if num[len(num)-1:len(num)] == "." {
		num = num[:len(num)-1]
	}
	return num

}
func getTime(timestring string) (string, error) {
	var timenum string
	if len(timestring) == 0 {
		t := time.Now()
		t1 := t.Format("20060102150405")
		year := t1[0:4]
		month := t1[4:6]
		day := t1[6:8]
		hour := t1[8:10]
		min := t1[10:12]
		second := t1[12:14]
		// 9/25/2019 14:05:01.732
		timenum = month + "/" + day + "/" + year + " " + hour + ":" + min + ":" + second
	} else {
		var timesnum, err = strconv.ParseInt(timestring[:len(timestring)-3], 10, 64)
		if err == nil {
			tm := time.Unix(timesnum, 0)
			timenum = tm.Format("01/02/2006 15:05:05")
		} else {
			timenum = timestring
		}
	}
	return timenum, nil
}
