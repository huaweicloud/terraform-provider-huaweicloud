package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResizeClusterRequestBodyExtendParam struct {

	// 专属云CCE集群可指定控制节点的规格
	DecMasterFlavor *string `json:"decMasterFlavor,omitempty"`

	// 是否自动扣款 - “true”：自动扣款 - “false”：不自动扣款 > 包周期集群时生效，不填写此参数时默认不会自动扣款。
	IsAutoPay *string `json:"isAutoPay,omitempty"`
}

func (o ResizeClusterRequestBodyExtendParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeClusterRequestBodyExtendParam struct{}"
	}

	return strings.Join([]string{"ResizeClusterRequestBodyExtendParam", string(data)}, " ")
}
