/*
 * @Version: 1.0
 * @Date: 2023-03-19 00:16:20
 * @LastEditTime: 2023-03-19 14:04:46
 */
package model

type DataOracles struct {
	QueryParam *DataOraclesQueryParam `json:"queryParam"` // 查询参数
	ResData    *ResData               `json:"resData"`    // 数据
}

type DataOraclesQueryParam struct {
	QueryStartTime string `json:"queryStartTime"` // 定位开始时间
	QueryEndTime   string `json:"queryEndTime"`   // 定位结束时间
}
