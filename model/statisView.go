/*
 * @Version: 1.0
 * @Date: 2023-03-19 11:11:30
 * @LastEditTime: 2023-03-19 14:02:36
 */
package model

type StatisViews struct {
	QueryParam *StatisViewsQueryParam `json:"queryParam"` // 查询参数
	ResData    *ResData               `json:"resData"`    // 数据
}

type StatisViewsQueryParam struct {
	QueryCycleType int    `json:"queryCycleType"` // 查询的周期类型 1天，2周，3月，4年
	QueryTime      string `json:"queryTime"`      // 传入的查询时间
	QueryStartTime string `json:"queryStartTime"` // 定位开始时间
	QueryEndTime   string `json:"queryEndTime"`   // 定位结束时间
}
