package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 数据信息详情
type Data struct {
	// 返回值类型。

	ResultType *string `json:"resultType,omitempty"`
	// 数据信息。

	Result *[]string `json:"result,omitempty"`
}

func (o Data) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Data struct{}"
	}

	return strings.Join([]string{"Data", string(data)}, " ")
}
