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
}

func (o TopReferSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopReferSummary struct{}"
	}

	return strings.Join([]string{"TopReferSummary", string(data)}, " ")
}
