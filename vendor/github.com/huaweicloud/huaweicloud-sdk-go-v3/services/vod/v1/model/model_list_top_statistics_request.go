package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTopStatisticsRequest struct {

	// 查询域名，暂只支持查询单个或者全部域名。  取值如下： - 单个加速域名，格式：example.test1.com。 - ALL：表示查询名下全部域名。
	Domain string `json:"domain"`

	// 查询日期，格式为yyyymmdd - date必须为昨天或之前的日期。 - 最多只能查最近一个月内的数据。
	Date string `json:"date"`
}

func (o ListTopStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTopStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListTopStatisticsRequest", string(data)}, " ")
}
