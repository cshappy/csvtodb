package svr_csvtosqlite

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	_ "github.com/mattn/go-sqlite3"
)

func Csvtosqlite(r *ghttp.Request) (interface{}, error) {
	var tablename = "test1"
	var Buf bytes.Buffer
	var result = []interface{}{}
	var result1 = []interface{}{}
	var dbfile = "C:\\Users\\m1894\\Desktop\\git\\test\\test.db"
	var Value = make(map[string]interface{})
	var i1 = 0
	var nameexist string
	var datasourceexist string
	var addressexist string
	// 根据字段名获取表单文件
	formFile, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return result, nil
	}

	defer formFile.Close()
	io.Copy(&Buf, formFile)
	contents := Buf.String()
	// fmt.Println(contents)
	arr := strings.Split(contents, "\n")

	db, _ := sql.Open("sqlite3", dbfile)
	rows, _ := db.Query("PRAGMA table_info([" + tablename + "])")
	var colum = []interface{}{}
	var columstring = ""
	for rows.Next() {
		var cid string
		var name string
		var type1 string
		var notnull string
		var dflt_value string
		var pk string
		rows.Scan(&cid, &name, &type1, &notnull, &dflt_value, &pk)
		colum = append(colum, []interface{}{name}...)
		columstring = columstring + "," + name
	}

	for i := 0; i < len(colum); i++ {
		if strings.ToUpper(colum[i].(string)) == "NAME" {
			nameexist = "exist"
		}
		if strings.ToUpper(colum[i].(string)) == "DATASOURCE" {
			datasourceexist = "exist"
		}
		if strings.ToUpper(colum[i].(string)) == "ADDRESS" {
			addressexist = "exist"
		}
	}
	// fmt.Println(nameexist, datasourceexist, addressexist)
	if nameexist != "exist" || datasourceexist != "exist" || addressexist != "exist" {
		result = append(result, "the important colum is not exist")
		// fmt.Println(result)
		return result, nil
	}

	var questionmark = "?"
	for i := 0; i < len(colum)-1; i++ {
		questionmark = questionmark + ",?"
	}
	//插入数据
	stmt, _ := db.Prepare("INSERT INTO " + tablename + "(" + columstring[1:] + ") values(" + questionmark + ")")
	tx, _ := db.Begin()
	var num = make(map[string]interface{})
	for i := 0; i < len(arr)-1; i++ {
		// fmt.Println(arr[i])
		record := cutcsv(arr[i])
		// record := strings.Split(arr[i], ",")
		if i1 == 0 {
			length := len(record)
			for x := 0; x < length; x++ {
				if record[x][0:1] == "#" {
					i1 = 0
					break
				}
				for y := 0; y < len(colum); y++ {
					if strings.ToUpper(colum[y].(string)) == strings.ToUpper(record[x]) {
						num[strings.ToUpper(record[x])] = x
						break
					}
					// if y == len(colum)-1 {
					// 	if strings.ToUpper(colum[y].(string)) == strings.ToUpper(record[x][:len(record[x])-1]) {
					// 		num[strings.ToUpper(record[x][:len(record[x])-1])] = y
					// 		fmt.Println(strings.ToUpper(record[x][:len(record[x])-1]))
					// 	}
					// }
				}

			}
			i1 = i1 + 1
		}
		fmt.Println(num)

		for a := 0; a < len(colum); a++ {
			number := num[strings.ToUpper(colum[a].(string))]
			// fmt.Println(record[number.(int)])
			if number == nil {
				Value[colum[a].(string)] = ""
			} else {
				Value[colum[a].(string)] = record[number.(int)]
			}

		}
		// fmt.Println(Value)
		if i == 0 {
			continue
		}
		for a := 0; a < len(colum); a++ {
			insertvalue, _ := Value[colum[a].(string)]
			result1 = append(result1, insertvalue)
		}
		_, _ = tx.Stmt(stmt).Exec(result1...)

		Value = make(map[string]interface{})
		result1 = []interface{}{}
	}
	tx.Commit()
	return result, nil
}

func cutcsv(strRemain string) []string {
	var column string
	var result []string
	for len(strRemain) > 0 {
		if strings.Count(strRemain, `,`) == 0 {
			result = append(result, strRemain)
			break
		} else {
			if strRemain[0:1] == `"` { //if starts with `"`
				strRemain = strRemain[1:] //remove first "
				column = ""
				for {
					cutEnd := strings.Index(strRemain, `",`) //it must end with `",`
					if cutEnd == -1 {                        //if there is no '",', it must be the last column and end with `"`.
						column = strings.TrimRight(strRemain, `"`)
						strRemain = ""
						break
					}
					column = column + strRemain[0:cutEnd] //assume this part is a complete column and does not include `",`
					strRemain = strRemain[cutEnd+2:]
					if strings.Count(column, `"`)%2 == 0 { // `"` must be pairs
						column = strings.Replace(column, `""`, `"`, len(column))
						break
					} else {
						column = column + `",` // if `"` are not in pairs, this `",` is not end and add it back to column since it is part of text.
					}
				}
			} else {
				cutEnd := strings.Index(strRemain, `,`)
				if cutEnd == -1 {
					cutEnd = len(strRemain)
				}
				column = strRemain[:cutEnd]
				strRemain = strRemain[cutEnd+1:]
			}
			result = append(result, column)
		}

	}
	return result
}
