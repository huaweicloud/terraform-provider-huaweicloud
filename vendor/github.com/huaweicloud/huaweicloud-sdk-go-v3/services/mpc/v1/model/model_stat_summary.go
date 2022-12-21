package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StatSummary struct {

	// 精确到小数点后两位。
	Value *float32 `json:"value,omitempty"`

	// 日期,精确到天,格式样例：2018/03/01。
	Date *string `json:"date,omitempty"`
}

func (o StatSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StatSummary struct{}"
	}

	return strings.Join([]string{"StatSummary", string(data)}, " ")
}
