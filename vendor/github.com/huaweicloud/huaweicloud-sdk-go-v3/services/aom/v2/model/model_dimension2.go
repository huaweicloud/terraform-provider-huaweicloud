package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指标维度参数
type Dimension2 struct {

	// 维度名称。
	Name string `json:"name"`

	// 指标数据的值。
	Value string `json:"value"`
}

func (o Dimension2) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Dimension2 struct{}"
	}

	return strings.Join([]string{"Dimension2", string(data)}, " ")
}
