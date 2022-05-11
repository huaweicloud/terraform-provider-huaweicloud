package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指标维度。
type Dimension struct {

	// 维度名称。
	Name string `json:"name"`

	// 维度取值。
	Value string `json:"value"`
}

func (o Dimension) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Dimension struct{}"
	}

	return strings.Join([]string{"Dimension", string(data)}, " ")
}
