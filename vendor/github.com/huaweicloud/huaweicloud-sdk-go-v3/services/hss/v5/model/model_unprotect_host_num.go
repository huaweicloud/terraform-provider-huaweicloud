package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 未防护服务器数
type UnprotectHostNum struct {
}

func (o UnprotectHostNum) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnprotectHostNum struct{}"
	}

	return strings.Join([]string{"UnprotectHostNum", string(data)}, " ")
}
