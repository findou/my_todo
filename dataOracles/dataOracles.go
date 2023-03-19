/*
 * @Version: 1.0
 * @Date: 2023-03-19 12:38:25
 * @LastEditTime: 2023-03-19 19:59:46
 */
package dataoracles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mytodo/data"
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"mytodo/utils"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DataOraclesHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodPost {
		post(w, r)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 上传文件,导入新的任务todo
func post(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// "text/plain; charset=utf-8"
	fType := http.DetectContentType(fileBytes)
	if fType != "text/plain; charset=utf-8" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("fType: %v\n", fType)
	dataOracle := &model.DataOracles{}
	err = json.Unmarshal(fileBytes, dataOracle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("dataOracle: %v\n", dataOracle)
	err = data.ImportData(dataOracle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func get(w http.ResponseWriter, r *http.Request) {
	var err error
	s := r.URL.EscapedPath()
	fmt.Println(s)
	uri := r.URL.RequestURI()
	urii, err := url.PathUnescape(uri)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uris := strings.Split(urii, "?")
	uriArr := strings.Split(uris[1], "&&")

	var startTime time.Time
	var startTimeStr string
	stStrArr := strings.Split(uriArr[0], "=")
	if stStrArr[0] == "startTime" {
		startTimeStr = stStrArr[1]
		startTime, err = utils.ParseStrTime(startTimeStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var endTime time.Time
	var endTimeStr string
	etStrArr := strings.Split(uriArr[1], "=")
	if etStrArr[0] == "endTime" {
		endTimeStr = etStrArr[1]
		endTime, err = utils.ParseStrTime(endTimeStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	fmt.Printf("startTime: %v\n", startTime)
	fmt.Printf("endTime: %v\n", endTime)

	jsonFile := false
	if len(uriArr) > 2 {
		//fileType=json
		ftStrArr := strings.Split(uriArr[2], "=")
		if ftStrArr[0] == "fileType" {
			typestr := ftStrArr[1]
			if typestr == "json" {
				jsonFile = true
			}
		}

	}

	queryParam := &model.DataOraclesQueryParam{
		QueryStartTime: utils.ParseTimeToStr(startTime),
		QueryEndTime:   utils.ParseTimeToStr(endTime),
	}

	resData, err := handletodo.HandleDataWithRangeTime(startTime, endTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := model.DataOracles{
		QueryParam: queryParam,
		ResData:    resData,
	}
	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if jsonFile {
		fileName := "任务完成情况数据-" + startTimeStr + "-" + endTimeStr + ".json"
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Disposition", "attachment;filename="+fileName)
	}
	w.Write(res)
	w.WriteHeader(http.StatusOK)
}
