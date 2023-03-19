/*
 * @Version: 1.0
 * @Date: 2023-03-18 18:12:17
 * @LastEditTime: 2023-03-19 20:06:31
 */
package statisViews

import (
	"encoding/json"
	handletodo "mytodo/handleTodo"
	"mytodo/model"
	"mytodo/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 周期内待办任务总览,包括待办完成情况
func CycleviewsHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func get(w http.ResponseWriter, r *http.Request) {
	var err error
	arr := strings.Split(r.URL.EscapedPath(), "/")
	queryCt := 1 //默认查询周期类型为 1 天，查询一天的数据
	cycleType := arr[2]
	if cycleType != "" {
		queryCt, err = strconv.Atoi(cycleType)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if queryCt >= 4 {
		queryCt = 1
	}
	queryTime := time.Now()
	if len(arr) > 3 {
		date := arr[3]
		if date != "" {
			// 将date转化为时间格式
			queryTime, err = time.ParseInLocation("2006-01-02", date, time.Now().Location())
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	st, et, err := handletodo.GetTimeArrangeByCt(queryTime, queryCt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	queryParam := &model.StatisViewsQueryParam{
		QueryCycleType: queryCt,
		QueryTime:      utils.ParseTimeToShortStr(queryTime),
		QueryStartTime: utils.ParseTimeToStr(st),
		QueryEndTime:   utils.ParseTimeToStr(et),
	}

	resData, err := handletodo.HandleDataWithRangeTime(st, et)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := &model.StatisViews{
		QueryParam: queryParam,
		ResData:    resData,
	}
	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
