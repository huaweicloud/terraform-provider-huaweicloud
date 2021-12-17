package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 克隆服务器类
type CloneServer struct {
	// 克隆服务器ID

	VmId *string `json:"vm_id,omitempty"`
	// 克隆虚拟机的名称

	Name *string `json:"name,omitempty"`
	// 克隆错误信息

	CloneError *string `json:"clone_error,omitempty"`
	// 克隆状态

	CloneState *string `json:"clone_state,omitempty"`
	// 克隆错误信息描述

	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o CloneServer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CloneServer struct{}"
	}

	return strings.Join([]string{"CloneServer", string(data)}, " ")
}
