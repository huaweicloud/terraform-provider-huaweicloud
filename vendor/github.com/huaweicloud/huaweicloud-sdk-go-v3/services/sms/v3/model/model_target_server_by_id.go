package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 目的端
type TargetServerById struct {
	// 目的端服务器ID

	VmId *string `json:"vm_id,omitempty"`
	// 目的端服务器名称

	Name *string `json:"name,omitempty"`
}

func (o TargetServerById) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetServerById struct{}"
	}

	return strings.Join([]string{"TargetServerById", string(data)}, " ")
}
