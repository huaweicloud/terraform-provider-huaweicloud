package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 克隆服务器基本信息
type CloneServerBrief struct {
	// 克隆服务器ID

	VmId *string `json:"vm_id,omitempty"`
	// 克隆虚拟机的名称

	Name *string `json:"name,omitempty"`
}

func (o CloneServerBrief) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloneServerBrief struct{}"
	}

	return strings.Join([]string{"CloneServerBrief", string(data)}, " ")
}
