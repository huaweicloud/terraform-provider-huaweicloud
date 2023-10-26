package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeCoreRequest Request Object
type UpgradeCoreRequest struct {

	// 指定待升级的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定待升级的节点类型，当前仅支持all。
	InstType string `json:"inst_type"`

	Body *UpgradingTheKernelBody `json:"body,omitempty"`
}

func (o UpgradeCoreRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeCoreRequest struct{}"
	}

	return strings.Join([]string{"UpgradeCoreRequest", string(data)}, " ")
}
