package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeNicSpec 节点网卡的描述信息。
type NodeNicSpec struct {
	PrimaryNic *NicSpec `json:"primaryNic,omitempty"`

	// 扩展网卡 >创建节点池添加节点时不支持该参数。
	ExtNics *[]NicSpec `json:"extNics,omitempty"`
}

func (o NodeNicSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeNicSpec struct{}"
	}

	return strings.Join([]string{"NodeNicSpec", string(data)}, " ")
}
