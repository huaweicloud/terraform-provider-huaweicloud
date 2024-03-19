package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TopReferSummary TOP100 Referer数据明细
type TopReferSummary struct {

	// referer值。
	Refer *string `json:"refer,omitempty"`

	// 对应查询类型的值。（流量单位：Byte）
	Value *int64 `json:"value,omitempty"`

	// 该referer的流量(或请求数)占当前查询条件下总流量(或请求数)的比例。保留4位小数
	Ratio *float64 `json:"ratio,omitempty"`
}

func (o TopReferSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopReferSummary struct{}"
	}

	return strings.Join([]string{"TopReferSummary", string(data)}, " ")
}
