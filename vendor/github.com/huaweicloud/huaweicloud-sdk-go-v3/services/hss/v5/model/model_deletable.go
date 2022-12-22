package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 是否允许删除该策略组
type Deletable struct {
}

func (o Deletable) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Deletable struct{}"
	}

	return strings.Join([]string{"Deletable", string(data)}, " ")
}
