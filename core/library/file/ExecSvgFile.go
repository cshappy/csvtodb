package lib_file

import (
	"archive/zip"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/beevik/etree"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gxml"
)

func ClearNameSpace(source string, destination string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(source); err != nil {
		panic(err)
	}
	// 删除源文件
	os.Remove("./tmp/" + source)
	root := doc.SelectElement("svg")
	searchChild(root)
	file, _ := os.Create(destination)
	doc.WriteTo(file)
}

// 遍历子元素 删除携带命名空间的元素及属性
func searchChild(root *etree.Element) {
	// 删除带gef的属性
	// 必须进行一次深拷贝， root.Attr是切片类型，每删除一次，其长度都会变化
	var attrArr = make([]etree.Attr, 0)
	for _, v := range root.Attr {
		attrArr = append(attrArr, v)
	}
	for i := 0; i < len(attrArr); i++ {
		if strings.Contains(attrArr[i].FullKey(), "gef:") {
			root.RemoveAttr(attrArr[i].FullKey())
		}
	}
	for _, v := range root.ChildElements() {
		if v.Space == "gef" {
			root.RemoveChildAt(v.Index())
		} else {
			if len(v.ChildElements()) > 0 {
				searchChild(v)
			}
		}
	}
}

// 文件压缩 file 源文件   destination 目的文件(压缩后的文件)
func Compress(file string, destination string) {
	d, _ := os.Create(destination)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	// 开始压缩
	f, _ := os.Open(file)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	header, _ := zip.FileInfoHeader(info)
	header.Name = "/" + header.Name
	writer, err := w.CreateHeader(header)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(writer, f)
	f.Close()
}

//获取json数据 content表示xml字节文件
func GetJSON(content []byte) []byte {
	var result = `{
		"result": []
	}`
	resultGjson, _ := gjson.DecodeToJson(result)
	var bindKeyValueRelation = getBindMap(content)
	var bindValueGjsonRealtion = getValueBindgJson(content)
	for bindkey, bindvalue := range bindKeyValueRelation {
		var obj = map[string]interface{}{
			"id":          bindkey,
			"description": nil,
		}
		// value 为最bindValue
		for _, value := range bindvalue {
			var Json = bindValueGjsonRealtion[value]
			if Json != nil {
				var animationProperty = Json.GetString("animatedProperty")
				switch animationProperty {
				case "foregroundColor", "backgroundColor", "edgeColor", "visible":
					obj[animationProperty] = map[string]interface{}{
						"dataSource":         Json.GetString("expression.math.semantics.apply.ci"),
						"errorMode":          Json.GetString("-errorMode"),
						"globalSharedTable":  Json.GetBool("-globalSharedTable"),
						"globalToggle":       Json.GetBool("-globalToggle"),
						"globalToggleSource": Json.GetString("-globalToggleSource"),
						"toggleRate":         Json.GetFloat64("-toggleRate"),
						"inputTolerance":     Json.GetString("expression.-tolerance"),
						"outputTolerance":    Json.GetString("-tolerance"),
					}
					// go语言不支持三元运算符
					var matchType string
					if Json.GetBool("level.-exactMatch") {
						matchType = "exactMatch"
					} else {
						matchType = "ranges"
					}
					// 必须使用反射，因为obj[animationProperty] 是interface类型
					obj[animationProperty].(map[string]interface{})["threshold"] = map[string]interface{}{
						"matchType":       matchType,
						"inputType":       "string",
						"outputType":      new(interface{}),
						"outOfRangeValue": Json.GetString("outOfRangeValue.#text"),
						"varyTable":       make(map[string]interface{}),
					}
					obj[animationProperty].(map[string]interface{})["threshold"].(map[string]interface{})["varyTable"] = getVaryTable(Json, obj[animationProperty].(map[string]interface{})["threshold"].(map[string]interface{}))
					break
				case "caption":
					obj[animationProperty] = map[string]interface{}{
						"dataSource":    Json.GetString("expression.math.semantics.apply.ci"),
						"errorMode":     Json.GetString("-errorMode"),
						"formatType":    Json.GetString("-formatType"),
						"justcfication": Json.GetString("-justification"),
						"interDigits":   Json.GetString("-integerDigits"),
						"decimalDigits": Json.GetString("-decimalDigits"),
						"dataEntryType": Json.GetString("-dataEntryType"),
						"confirmation":  Json.GetBool("-confirmation"),
					}
					break
				case "horizontalPosition", "verticalPosition":
					obj[animationProperty] = map[string]interface{}{
						"dataSource":     Json.GetString("expression.math.semantics.apply.ci"),
						"motionType":     "absolute",
						"errorMode":      Json.GetString("-errorMode"),
						"position":       Json.GetString("-horizontalPosition"),
						"lowInputValue":  Json.GetString("-lowInputValue"),
						"highInputValue": Json.GetString("-highInputValue"),
						"minOutputValue": Json.GetString("-minOutputValue"),
						"maxOutputValue": Json.GetString("-maxOutputValue"),
						"autoFetchInput": Json.GetBool("-autoFetchInput"),
					}
					break
				}
			}
		}
		resultGjson.Append("result", obj)
	}
	bytes, _ := resultGjson.ToJson()
	return bytes
}

// 生成varytable 字段
func getVaryTable(g *gjson.Json, obj map[string]interface{}) []map[string]interface{} {
	var level = g.GetInterfaces("level")

	var result = make([]map[string]interface{}, 0)
	for _, v := range level {

		var levelMap = v.(map[string]interface{})
		var outputType = levelMap["outputValue"].(map[string]interface{})["-type"].(string)
		var toggleValue string
		if levelMap["toggleValue"] == nil {
			toggleValue = ""
		} else {
			toggleValue = levelMap["toggleValue"].(map[string]interface{})["#text"].(string)
		}
		var outputvalue interface{}
		if outputType == "string" {
			outputvalue = g.GetBool("outOfRangeValue.#text")
			obj["outputType"] = "bool"
		} else {
			obj["outputType"] = "color"
			outputvalue = map[string]interface{}{
				"color":   levelMap["outputValue"].(map[string]interface{})["#text"],
				"opacity": 1,
			}
		}
		var temp = map[string]interface{}{
			"loMatchValue": levelMap["loMatchValue"].(map[string]interface{})["#text"],
			"outputvalue":  outputvalue,
			"toggleOutput": true,
			"toggleValue": map[string]interface{}{
				"color":   toggleValue,
				"opacity": 1,
			},
		}
		result = append(result, temp)
	}
	return result
}

// 获取绑定关系 bindkey => bindvalue
func getBindMap(content []byte) map[string][]string {
	var bindMap = make(map[string][]string)
	result, err := gxml.Decode(content)
	if err != nil {
		panic(err.Error())
	}
	var foreignObject = result["svg"].(map[string]interface{})["foreignObject"]
	var bindings = foreignObject.(map[string]interface{})["bindings"]
	var binding = bindings.(map[string]interface{})["binding"].([]interface{})
	for _, v := range binding {
		var row = v.(map[string]interface{})
		var bindKey = row["bindKey"].(string)
		if bindMap[bindKey] != nil {
			bindMap[bindKey] = append(bindMap[bindKey], row["bindValue"].(string))
		} else {
			var bindValue = make([]string, 0)
			switch reflect.TypeOf(row["bindValue"]).String() {
			case "string":
				bindValue = append(bindValue, row["bindValue"].(string))
				break
			case "[]interface {}":
				for _, v3 := range row["bindValue"].([]interface{}) {
					bindValue = append(bindValue, v3.(string))
				}
				break
			}
			bindMap[bindKey] = bindValue
		}
	}
	return bindMap
}

// 返回bindValue 对应的gjson引用
func getValueBindgJson(content []byte) map[string]*gjson.Json {
	var result = make(map[string]*gjson.Json)
	bytes, _ := gxml.ToJson(content)
	g, _ := gjson.DecodeToJson(bytes)
	var gStrArr = g.GetStrings("svg.foreignObject.auxiliaryObjects.animations.lookupAnimation")
	for _, v := range gStrArr {
		g2, _ := gjson.DecodeToJson(v)
		// expression.math.semantics.apply.ci
		var id = g2.GetString("-id")
		result[id] = g2
	}
	return result
}
